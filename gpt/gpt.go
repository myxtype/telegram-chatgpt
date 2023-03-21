package gpt

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/cast"
	"io"
	"log"
	"net/http"
	"telegram-chatgpt/conf"
)

const BASEURL = "https://api.openai.com"

type ChatGPTResponseError struct {
	Message string `json:"message"`
}

type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
	Error   *ChatGPTResponseError    `json:"error,omitempty"`
}

type ChatGPTRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
	User             string  `json:"user"`

	Messages []*ChatGPTRequestMessage `json:"messages,omitempty"`
	Prompt   string                   `json:"prompt,omitempty"`
}

func Completions(user int64, msg string) (string, error) {
	requestBody := ChatGPTRequestBody{
		Model:            "gpt-3.5-turbo",
		Messages:         GetSessionMessages(user, msg),
		MaxTokens:        conf.Config().ChatGPT.MaxTokens,
		Temperature:      conf.Config().ChatGPT.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		User:             cast.ToString(user),
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"/v1/completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := conf.Config().ChatGPT.ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}

	if gptResponseBody.Error != nil {
		return gptResponseBody.Error.Message, nil
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v["text"].(string)
			break
		}
	}

	if reply != "" {
		SaveSessionMessage(user, msg, reply)
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}
