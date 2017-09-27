package session

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_NewMgr(t *testing.T) {
	request := &http.NewRequest{
		method: "/",
		url:    "/",
		body:   strings.NewReader(""),
	}
	mgr := NewSessionMgr("123ccc", 1)
	mgr.Set("name", "ccc", "234")
	val, ok := mgr.Get("name", "ccc")
	if !ok || val.(string) != "234" {
		t.Errorf("get things error")
	}
	time.Sleep(time.Duration(1+1) * time.Second)
	if _, ok = mgr.Get("name", "ccc"); ok {
		t.Errorf("gc failed")
	}
}
