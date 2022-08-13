package impl

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramBotAPI struct {
	channelID int64
	adminID   int64
	bot       *tgbotapi.BotAPI
}

func NewTelegramBot(config TelegramConfig) (*telegramBotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		return nil, err
	}

	return &telegramBotAPI{
		channelID: config.ChannelID,
		adminID:   config.AdminID,
		bot:       bot,
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

func (tg *telegramBotAPI) SendMessageToAdmin(text string) error {
	return tg.SendMessage(tg.adminID, text)
}

func (tg *telegramBotAPI) SendPostsToChannel(ctx context.Context, posts []Post) {
	for _, post := range posts {
		txt := fmt.Sprintf("%s\n<i>%s</i>\n<b>%s</b>", post.Link, strings.Join(post.Tags, ", "), post.PublishedAt)
		err := tg.SendMessage(tg.channelID, txt)
		if err != nil {
			log.Println("SendPostsToChannel error:", err)
		}
	}
}
