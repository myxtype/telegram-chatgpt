package gpt

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"telegram-chatgpt/conf"
	"telegram-chatgpt/session"
)

const BASEURL = "https://api.openai.com"

type ChatGPTResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}

type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
	User             string  `json:"user"`
}

func Completions(user, msg string) (string, error) {
	requestBody := ChatGPTRequestBody{
		Model:            "text-davinci-003",
		Prompt:           session.GetPrompt(user, msg),
		MaxTokens:        conf.Config().ChatGPT.MaxTokens,
		Temperature:      conf.Config().ChatGPT.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		User:             user,
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

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v["text"].(string)
			break
		}
	}

	if reply != "" {
		session.SaveMsg(user, msg, reply)
	}
	log.Printf("gpt response text: %s \n", reply)
	return reply, nil
}