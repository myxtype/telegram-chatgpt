package conf

import "time"

type GbeConfig struct {
	Bot     botConfig
	Session sessionConfig
	ChatGPT chatGPTConfig
	Limiter limiterConfig
}

type botConfig struct {
	Token string

	HelloText        string
	SessionClearText string
	LimiterText      string
	ThinkingText     string
}

type sessionConfig struct {
	TokensLimit int
	Exp         time.Duration
}

type chatGPTConfig struct {
	ApiKey      string
	MaxTokens   int
	Temperature float32
}

type limiterConfig struct {
	Tokens   uint64
	Interval time.Duration
}
