package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMsg(botAPI *tgbotapi.BotAPI, chat_id int64, text string) error {
	msg := tgbotapi.NewMessage(chat_id, text)
	_, err := botAPI.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
