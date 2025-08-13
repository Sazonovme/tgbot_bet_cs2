package handler

import (
	"RushBananaBet/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMsg(botAPI *tgbotapi.BotAPI, chat_id int64, text string, keyboard [][]tgbotapi.KeyboardButton) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chat_id, text)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	botMsg, err := botAPI.Send(msg)
	if err != nil {
		logger.Error("Err send msg", "helpers-sendMsg()", err)
		return botMsg, err
	}
	return botMsg, nil
}

func deleteMsg(botAPI *tgbotapi.BotAPI, chat_id int64, message_id int) error {

	deleteConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    chat_id,
		MessageID: message_id,
	}

	if _, err := botAPI.Send(deleteConfig); err != nil {
		logger.Error("Err delete msg", "helpers-deleteMsg()", err)
		return err
	}

	return nil
}

func deleteKeyboard(botAPI *tgbotapi.BotAPI, chat_id int64, message_id int) error {
	edit := tgbotapi.NewEditMessageReplyMarkup(
		chat_id,
		message_id,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
	)

	if _, err := botAPI.Send(edit); err != nil {
		logger.Error("Err delete keyboard", "helpers-deleteKeyboard()", err)
		return err
	}

	return nil

}
