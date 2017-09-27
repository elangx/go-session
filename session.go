package session

import (
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const KEY_LENGTH = 64

type SessionMgr struct {
	mSessions map[string]*Session
	mKeyName  string
	mLock     sync.RWMutex
	mMaxTime  int64
	store     Store
}

func NewCookieSessionMgr(keyName string, maxTime int64) *SessionMgr {
	mgr := &SessionMgr{
		mKeyName:  keyName,
		mMaxTime:  maxTime,
		mSessions: make(map[string]*Session),
		store:     CookieStore{KeyName: keyName, ExpTime: maxTime},
	}
	mgr.GC()
	return mgr
}

func NewHeaderSessionMgr(keyName string, maxTime int64) *SessionMgr {
	mgr := &SessionMgr{
		mKeyName:  keyName,
		mMaxTime:  maxTime,
		mSessions: make(map[string]*Session),
		store:     HeaderStore{KeyName: keyName},
	}
	mgr.GC()
	return mgr
}

func generCookieValue() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < KEY_LENGTH; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (p *SessionMgr) GC() {
	p.mLock.Lock()
	defer p.mLock.Unlock()

	for id, session := range p.mSessions {
		if session.mLastVisitTime.Unix() > p.mMaxTime&time.Now().Unix() {
			delete(p.mSessions, id)
		}
	}
	time.AfterFunc(time.Duration(p.mMaxTime)*time.Second, func() {
		p.GC()
	})
}

func (p *SessionMgr) Set(r *http.Request, w http.ResponseWriter, key, val interface{}) {
	storeKey := p.store.Get(r)
	if storeKey == "" {
		storeKey := generCookieValue()
		p.store.Set(w, storeKey)
	}

	p.mLock.Lock()
	defer p.mLock.Unlock()

	if _, ok := p.mSessions[storeKey]; !ok {
		p.mSessions[storeKey] = &Session{
			mLastVisitTime: time.Now(),
			Value:          make(map[interface{}]interface{}),
		}
	}
	p.mSessions[storeKey].Value[key] = val
}

func (p *SessionMgr) Get(r *http.Request, key interface{}) (interface{}, bool) {
	storeKey := p.store.Get(r)
	if storeKey == "" {
		return nil, false
	}

	p.mLock.Lock()
	defer p.mLock.Unlock()

	if session, ok := p.mSessions[storeKey]; ok {
		if val, ok := session.Value[key]; ok {
			session.mLastVisitTime = time.Now()
			return val, ok
		}
	}
	return nil, false
}

type Session struct {
	mLastVisitTime time.Time
	Value          map[interface{}]interface{}
}
