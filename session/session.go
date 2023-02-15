package session

import (
	"sync"
	"telegram-chatgpt/conf"
	"time"
	"unicode/utf8"
)

var sessions = sync.Map{}

type session struct {
	records  []*record
	lastTime time.Time
}

type record struct {
	q     string
	a     string
	count int
}

func newRecord(q, a string) *record {
	return &record{
		q:     q,
		a:     a,
		count: utf8.RuneCountInString(q) + utf8.RuneCountInString(a),
	}
}

func GetPrompt(user, msg string) string {
	var sm string

	sess := getSession(user)
	for _, r := range sess.records {
		sm += "Q:" + r.q + "\nA:" + r.a + "\n"
	}

	sm += "Q:" + msg + "\nA:"

	return sm
}

func SaveMsg(user, msg, reply string) {
	s := getSession(user)
	s.records = append(s.records, newRecord(msg, reply))

	var total int
	for _, n := range s.records {
		total += n.count
	}

	for _, n := range s.records {
		if total > conf.Config().Session.TokensLimit {
			total -= n.count
			s.records = s.records[1:]
		}
	}
}

func ClearSession(user string) {
	sessions.Delete(user)
}

func getSession(user string) *session {
	var sess *session

	v, fond := sessions.Load(user)
	if fond {
		sess = v.(*session)
	} else {
		sess = &session{
			records:  []*record{},
			lastTime: time.Now(),
		}
		sessions.Store(user, sess)
	}

	// exp
	if sess.lastTime.Before(time.Now().Add(-1 * conf.Config().Session.Exp * time.Minute)) {
		sess.records = []*record{}
	}

	return sess
}
