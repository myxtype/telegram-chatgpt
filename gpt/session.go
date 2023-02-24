package gpt

import (
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
	"telegram-chatgpt/conf"
	"time"
	"unicode/utf8"
)

var sessions = cache.New(5*time.Minute, 10*time.Minute)

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
	sessions.Delete(cast.ToString(user))
}

func getSession(user int64) *session {
	var sess *session

	v, fond := sessions.Get(cast.ToString(user))
	if fond {
		sess = v.(*session)
	} else {
		sess = &session{
			records:  []*record{},
			lastTime: time.Now(),
		}
		sessions.SetDefault(cast.ToString(user), sess)
	}

	if sess.lastTime.Before(time.Now().Add(-1 * conf.Config().Session.Exp * time.Minute)) {
		sess.records = []*record{}
	}

	return sess
}
