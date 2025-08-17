package handler

import (
	"RushBananaBet/internal/model"
	"RushBananaBet/internal/ui"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMsg(botAPI *tgbotapi.BotAPI, chat_id int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chat_id, text)
	if len(keyboard.InlineKeyboard) != 0 {
		msg.ReplyMarkup = keyboard
	}
	botMsg, err := botAPI.Send(msg)
	if err != nil {
		return botMsg, err
	}
	return botMsg, nil
}

func CalcPointsForBet(prediction string, result string) string {
	points := "0"
	if prediction == "1" {
		if strings.HasPrefix(result, "2") {
			points = "1"
		}
	} else if prediction == "2" {
		if strings.HasPrefix(result, "0") || strings.HasPrefix(result, "1") {
			points = "1"
		}
	} else if prediction == result {
		points = "2"
	}

	return points
}

func openMainMenu(botAPI *tgbotapi.BotAPI, messageText string, userData *model.User) error {
	if messageText != "" {
		_, err := sendMsg(botAPI, userData.Chat_id, messageText, tgbotapi.InlineKeyboardMarkup{})
		if err != nil {
			return err
		}
	}

	keyboard := ui.PaintMainMenu(userData.IsAdmin)
	msg, err := sendMsg(botAPI, userData.Chat_id, "Главное меню:", keyboard)
	if err != nil {
		return err
	}
	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, []int{msg.MessageID}, "main_menu")

	return nil
}
