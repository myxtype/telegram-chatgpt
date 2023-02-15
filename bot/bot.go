package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"telegram-chatgpt/conf"
	"telegram-chatgpt/gpt"
	"telegram-chatgpt/session"
)

func Start() {
	bot, err := tgbotapi.NewBotAPI(conf.Config().Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				if update.Message.Command() == "clear" {
					session.ClearSession(update.Message.From.UserName)
				}
				continue
			}

			var text string

			if update.Message.Chat.IsPrivate() {
				text = update.Message.Text
			} else if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
				text = update.Message.Text
			} else {
				atBotText := fmt.Sprintf("@%s", bot.Self.UserName)

				if strings.Index(update.Message.Text, atBotText) < 0 {
					continue
				}

				text = strings.ReplaceAll(update.Message.Text, atBotText, "")
			}

			reply, err := gpt.Completions(update.Message.From.UserName, strings.TrimSpace(text))
			if err != nil {
				log.Printf("%s", err.Error())
				continue
			}

			if reply != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
		}
	}
}
