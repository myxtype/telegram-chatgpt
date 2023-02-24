package gpt

import (
	"github.com/golang/groupcache/lru"
	"telegram-chatgpt/conf"
	"time"
	"unicode/utf8"
)

var sessions = lru.New(1024)

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

func GetPrompt(user int64, msg string) string {
	var sm string

	sess := getSession(user)
	for _, r := range sess.records {
		sm += "Q:" + r.q + "\nA:" + r.a + "\n"
	}

	sm += "Q:" + msg + "\nA:"

	return sm
}

func SaveMsg(user int64, msg, reply string) {
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

func GetSessionRecordsCount(user int64) int {
	sess := getSession(user)
	return len(sess.records)
}

func ClearSession(user int64) {
	sessions.Remove(user)
}

func getSession(user int64) *session {
	var sess *session

	v, fond := sessions.Get(user)
	if fond {
		sess = v.(*session)
	} else {
		sess = &session{
			records:  []*record{},
			lastTime: time.Now(),
		}
		sessions.Add(user, sess)
	}

	// exp
	if sess.lastTime.Before(time.Now().Add(-1 * conf.Config().Session.Exp * time.Minute)) {
		sess.records = []*record{}
	}

	return sess
}
