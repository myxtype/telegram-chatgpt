package conf

import "time"

type GbeConfig struct {
	Bot     botConfig
	Session sessionConfig
	ChatGPT chatGPTConfig
}

type botConfig struct {
	Token string
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
