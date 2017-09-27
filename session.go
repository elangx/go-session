package session

import (
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

func (p *SessionMgr) Set(id string, key, val interface{}) {
	p.mLock.Lock()
	defer p.mLock.Unlock()

	if _, ok := p.mSessions[id]; !ok {
		p.mSessions[id] = &Session{
			mLastVisitTime: time.Now(),
			Value:          make(map[interface{}]interface{}),
		}
	}
	p.mSessions[id].Value[key] = val
}

func (p *SessionMgr) Get(id string, key interface{}) (interface{}, bool) {
	p.mLock.Lock()
	defer p.mLock.Unlock()

	if session, ok := p.mSessions[id]; ok {
		if val, ok := session.Value[key]; ok {
			return val, ok
		}
	}
	return nil, false
}

type Session struct {
	mLastVisitTime time.Time
	Value          map[interface{}]interface{}
}
