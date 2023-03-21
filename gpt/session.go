package gpt

import (
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
	"telegram-chatgpt/conf"
	"time"
	"unicode/utf8"
)

var sessions = cache.New(30*time.Minute, 10*time.Minute)

type Session struct {
	records  []*Record
	lastTime time.Time
}

type Record struct {
	q     string
	a     string
	count int
}

func NewRecord(q, a string) *Record {
	return &Record{
		q:     q,
		a:     a,
		count: utf8.RuneCountInString(q) + utf8.RuneCountInString(a),
	}
}

func GetSessionPrompt(user int64, msg string) string {
	var sm string

	sess := GetSession(user)
	for _, r := range sess.records {
		sm += "Q:" + r.q + "\nA:" + r.a + "\n"
	}

	sm += "Q:" + msg + "\nA:"

	return sm
}

func GetSessionMessages(user int64, msg string) []*ChatGPTMessage {
	var messages []*ChatGPTMessage

	if conf.Config().ChatGPT.Foreword != "" {
		messages = append(messages, &ChatGPTMessage{
			Role:    "system",
			Content: conf.Config().ChatGPT.Foreword,
		})
	}

	sess := GetSession(user)
	for _, r := range sess.records {
		messages = append(messages,
			&ChatGPTMessage{
				Role:    "user",
				Content: r.a,
			},
			&ChatGPTMessage{
				Role:    "assistant",
				Content: r.q,
			})
	}

	messages = append(messages, &ChatGPTMessage{
		Role:    "user",
		Content: msg,
	})

	return messages
}

func SaveSessionMessage(user int64, msg, reply string) {
	s := GetSession(user)
	s.records = append(s.records, NewRecord(msg, reply))

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
	sess := GetSession(user)
	return len(sess.records)
}

func ClearSession(user int64) {
	sessions.Delete(cast.ToString(user))
}

func GetSession(user int64) *Session {
	var sess *Session

	v, fond := sessions.Get(cast.ToString(user))
	if fond {
		sess = v.(*Session)
	} else {
		sess = &Session{
			records:  []*Record{},
			lastTime: time.Now(),
		}
		sessions.SetDefault(cast.ToString(user), sess)
	}

	if sess.lastTime.Before(time.Now().Add(-1 * conf.Config().Session.Exp * time.Minute)) {
		sess.records = []*Record{}
	}

	return sess
}
