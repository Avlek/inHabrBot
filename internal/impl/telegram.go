package impl

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramBotAPI struct {
	adminChatID int64
	bot         *tgbotapi.BotAPI
}

func NewTelegramBot(adminChatID int64, apiKey string) (*telegramBotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	return &telegramBotAPI{
		adminChatID,
		bot,
	}, nil
}

func (tg *telegramBotAPI) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := tg.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (tg *telegramBotAPI) SendMessageToAdmin(text string) error {
	msg := tgbotapi.NewMessage(tg.adminChatID, text)
	_, err := tg.bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return err
}
