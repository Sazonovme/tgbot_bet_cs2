package handler

import (
	"RushBananaBet/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PrepareUserData(update *tgbotapi.Update) *model.User {

	if update.CallbackQuery != nil {
		return &model.User{
			Chat_id:      update.CallbackQuery.Message.Chat.ID,
			Username:     update.CallbackQuery.Message.From.UserName,
			First_name:   update.CallbackQuery.Message.From.FirstName,
			Last_name:    update.CallbackQuery.Message.From.LastName,
			CallbackData: update.CallbackQuery.Data,
			IsAdmin:      model.IsAdmin(update.CallbackQuery.Message.From.UserName),
		}
	} else {
		return &model.User{
			Chat_id:    update.Message.Chat.ID,
			Username:   update.Message.From.UserName,
			First_name: update.Message.From.FirstName,
			Last_name:  update.Message.From.LastName,
			TextMsg:    update.Message.Text,
			IsAdmin:    model.IsAdmin(update.Message.From.UserName),
		}
	}

}
