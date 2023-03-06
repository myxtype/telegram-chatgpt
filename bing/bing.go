package bing

const baseUri = ""

type ChatBing struct {
	cookie string
}

func NewChatBing() *ChatBing {
	return &ChatBing{}
}
