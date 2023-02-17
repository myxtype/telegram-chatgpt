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
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "clear":
					session.ClearSession(update.Message.From.UserName)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, conf.Config().Bot.SessionClearText)
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				case "start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, conf.Config().Bot.HelloText)
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				}
				continue
			}

			// Limier take
			_, _, rest, ok, err := limiter.Take(context.Background(), update.Message.From.UserName)
			if err != nil {
				log.Printf("limiter error %s", err.Error())
				continue
			}

			if !ok {
				sub := time.UnixMicro(int64(rest / 1000)).Sub(time.Now())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					fmt.Sprintf(conf.Config().Bot.LimiterText, sub.Seconds()),
				)
				if !update.Message.Chat.IsPrivate() {
					msg.ReplyToMessageID = update.Message.MessageID
				}

				replyMsg, err := bot.Send(msg)
				if err == nil {
					time.AfterFunc(sub, func() {
						bot.Send(tgbotapi.NewDeleteMessage(replyMsg.Chat.ID, replyMsg.MessageID))
					})
				}

				continue
			}

			var text string

			if update.Message.Chat.IsPrivate() {
				text = update.Message.Text
			} else if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.UserName == bot.Self.UserName {
				text = update.Message.Text
			} else if strings.Index(update.Message.Text, "/") == 0 {
				text = strings.Replace(update.Message.Text, "/", "", 1)
			} else {
				atBotText := fmt.Sprintf("@%s", bot.Self.UserName)

				if strings.Index(update.Message.Text, atBotText) < 0 {
					continue
				}

				text = strings.ReplaceAll(update.Message.Text, atBotText, "")
			}

			// ReplyToMessage
			if update.Message.ReplyToMessage != nil {
				if session.GetSessionRecordsCount(update.Message.From.UserName) == 0 {
					text = update.Message.ReplyToMessage.Text + "\n" + text
				}
			}

			// Thinking
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, conf.Config().Bot.ThinkingText)
			if !update.Message.Chat.IsPrivate() {
				msg.ReplyToMessageID = update.Message.MessageID
			}

			thinkMsg, err := bot.Send(msg)
			if err != nil {
				continue
			}

			// call ChatGPT
			reply, err := gpt.Completions(update.Message.From.UserName, strings.TrimSpace(text))
			if err != nil {
				log.Printf("gpt completions error %s", err.Error())
				bot.Send(tgbotapi.NewEditMessageText(thinkMsg.Chat.ID, thinkMsg.MessageID, err.Error()))
				continue
			}

			if reply != "" {
				bot.Send(tgbotapi.NewEditMessageText(thinkMsg.Chat.ID, thinkMsg.MessageID, reply))
			} else {
				bot.Send(tgbotapi.NewEditMessageText(thinkMsg.Chat.ID, thinkMsg.MessageID, "没有得到任何消息！"))
			}
		}
	}
}
