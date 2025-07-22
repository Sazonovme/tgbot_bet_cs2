package handler

import (
	"RushBananaBet/internal/model"
	"errors"

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

func middlewareAuth(botAPI *tgbotapi.BotAPI, userData *model.User) error {
	if model.IsAdmin(userData.Username) {
		return nil
	}
	sendMsg(botAPI, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов")
	return errors.New("This user is not admin")
}
