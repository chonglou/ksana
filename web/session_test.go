package ksana_web

import (
	"testing"
)

const sid = "test_session_sid"

func TestFileSession(t *testing.T) {
	session_t(&FileSessionProvider{path: "/tmp/ksana/tmp/sessions"}, t)
}

func session_t(sp SessionProvider, t *testing.T) {
	sess, err := sp.Read(sid)
	if err != nil {
		t.Errorf("Session init error: %v", err)
	}
	key, val := "aaa", 1234
	err = sess.Set(key, val)
	if err != nil {
		t.Errorf("Session set error: %v", err)
	}

	s1, e1 := sp.Read(sid)
	if e1 != nil {
		t.Errorf("Session read error: %v", e1)
	}

	if s1.Get(key) != val {
		t.Errorf("Want %i, Get %i", val, s1.Get(key))
	}

}
