package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sethvargo/go-limiter/memorystore"
	"log"
	"strings"
	"telegram-chatgpt/conf"
	"telegram-chatgpt/gpt"
	"telegram-chatgpt/session"
	"time"
)

func Start() {
	bot, err := tgbotapi.NewBotAPI(conf.Config().Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	limiter, err := memorystore.New(&memorystore.Config{
		Tokens:   conf.Config().Limiter.Tokens,
		Interval: time.Minute * conf.Config().Limiter.Interval,
	})
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				if update.Message.Command() == "clear" {
					session.ClearSession(update.Message.From.UserName)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "session cleared")
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
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

			// limier take
			_, _, rest, ok, err := limiter.Take(context.Background(), update.Message.From.UserName)
			if err != nil {
				log.Printf("limiter error %s", err.Error())
				continue
			}

			if !ok {
				sub := time.UnixMicro(int64(rest / 1000)).Sub(time.Now())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					fmt.Sprintf("限制每%d分钟%d次请求(剩余%v秒)", conf.Config().Limiter.Interval, conf.Config().Limiter.Tokens, sub.Seconds()),
				)
				msg.ReplyToMessageID = update.Message.MessageID

				// clear waring msg
				rmsg, err := bot.Send(msg)
				if err == nil {
					time.AfterFunc(sub, func() {
						bot.Send(tgbotapi.NewDeleteMessage(rmsg.Chat.ID, rmsg.MessageID))
					})
				}

				continue
			}

			// call ChatGPT
			reply, err := gpt.Completions(update.Message.From.UserName, strings.TrimSpace(text))
			if err != nil {
				log.Printf("gpt completions error %s", err.Error())
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
