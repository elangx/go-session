package session

import (
	"net/http"
	"sync"
	"time"
)

type SessionMgr struct {
	mSessions   map[string]*Session
	mCookieName string
	mLock       sync.RWMutex
	mMaxTime    int64
}

func NewSessionMgr(cookieName string, maxTime int64) *SessionMgr {
	mgr := &SessionMgr{
		mCookieName: cookieName,
		mMaxTime:    maxTime,
		mSessions:   make(map[string]*Session),
	}
	mgr.GC()
	return mgr
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

func (p *SessionMgr) Set(w http.ResponseWriter, key, val interface{}) {
	cookie := &http.Cookie{
		Name:    key.(string),
		Value:   val.(string),
		Expires: time.Now().Add(maxTime.(time.Duration)),
	}

	p.mLock.Lock()
	defer p.mLock.Unlock()

	if _, ok := p.mSessions[cookie.Value]; !ok {
		p.mSessions[id] = &Session{
			mLastVisitTime: time.Now(),
			Value:          make(map[interface{}]interface{}),
		}
	}
	p.mSessions[id].Value[key] = val
	w.SetCookie(w, cookie)
}

func (p *SessionMgr) Get(r *http.Request, key interface{}) (interface{}, bool) {
	cookie, err := r.Cookie(p.mCookieName)
	if err != nil {
		return nil, false
	}

	p.mLock.Lock()
	defer p.mLock.Unlock()

	if session, ok := p.mSessions[cookie.Value]; ok {
		if val, ok := session.Value[key]; ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *SessionMgr) Del(r *http.Request, w http.ResponseWriter) {
	cookie, err := r.Cookie(p.mCookieName)
	if err != nil {
		return
	}

	p.mLock.Lock()
	defer p.mLock.Unlock()

	if session, ok := p.mSessions[cookie.Value]; ok {
		delete(p.mSessions, cookie.Value)
	}
}

type Session struct {
	mLastVisitTime time.Time
	Value          map[interface{}]interface{}
}
