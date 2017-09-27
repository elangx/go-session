package session

import (
	"net/http"
	"time"
)

type Store interface {
	Get(r *http.Request) string
	Set(w http.ResponseWriter)
}

type CookieStore struct {
	ExpTime int64
	KeyName string
}

func (p *CookieStore) Get(r *http.Request) string {
	cookie, err := r.Cookie(p.KeyName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (p *CookieStore) Set(w http.Writer, val string) {
	cookie := &http.Cookie{
		Name:    p.KeyName,
		Value:   val,
		Expires: time.Now().Add(time.Duration(p.ExpTime)),
	}
	http.SetCookie(w, cookie)
}

type HeaderStore struct {
	KeyName string
}

func (p *HeaderStore) Get(r *http.Request) string {
	return r.Header.Get(p.KeyName)
}

func (p *HeaderStore) Set(w http.Writer, val string) {
	w.Header().Set(p.KeyName, val)
}
