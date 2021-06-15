package impl

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramBotAPI struct {
	channelID int64
	bot       *tgbotapi.BotAPI
}

func NewTelegramBot(channelID int64, apiKey string) (*telegramBotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	return &telegramBotAPI{
		channelID,
		bot,
	}, nil
}

func (tg *telegramBotAPI) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := tg.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (tg *telegramBotAPI) SendPostsToChannel(ctx context.Context, posts []Post) {
	for _, post := range posts {
		err := tg.SendMessage(tg.channelID, fmt.Sprintf("%s\n<b>%s</b>", post.Link, post.PublishedAt))
		if err != nil {
			log.Println("SendPostsToChannel error:", err)
		}
	}
}
