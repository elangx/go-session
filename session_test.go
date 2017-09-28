package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_CookieMgr(t *testing.T) {
	cookie := &http.Cookie{
		Name:    "test-test",
		Value:   "test-cookie",
		Expires: time.Now().AddDate(0, 1, 0),
	}
	r, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Errorf(err.Error())
	}
	w := httptest.NewRecorder()
	r.AddCookie(cookie)
	mgr := NewCookieSessionMgr("test-test", 1)
	mgr.Set(r, w, "testkey", "testval")
	resultCookies := w.Result().Cookies()
	bRc := false
	for _, rc := range resultCookies {
		if rc.Name == "test-test" && rc.Value == "test-cookie" {
			bRc = true
		}
	}
	if !bRc {
		t.Errorf("get result cookie error")
	}
	val, ok := mgr.Get(r, "testkey")
	if !ok || val.(string) != "testval" {
		t.Errorf("get testkey error")
	}
	time.Sleep(time.Duration(1+1) * time.Second)
	if _, ok = mgr.Get(r, "testkey"); ok {
		t.Errorf("gc failed")
	}
}

func Test_HeaderMgr(t *testing.T) {
	r, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Errorf(err.Error())
	}
	r.Header.Add("test-test", "test-header")
	w := httptest.NewRecorder()
	mgr := NewHeaderSessionMgr("test-test", 1)
	mgr.Set(r, w, "testkey", "testval")
	val, ok := mgr.Get(r, "testkey")
	if !ok || val.(string) != "testval" {
		t.Errorf("get header test key error")
	}
	val = w.Header().Get("test-test")
	if val != "test-header" {
		t.Errorf("get header result error: %s", val)
	}
	time.Sleep(time.Duration(2) * time.Second)
	if _, ok = mgr.Get(r, "testkey"); ok {
		t.Errorf("gc failed")
	}
}
