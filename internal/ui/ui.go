package ui

import (
	"RushBananaBet/internal/model"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMainMenuMsg(chat_id int64, textMessage string, isAdmin bool) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, textMessage)
	msg.ReplyMarkup = GetMainMenuKeyboard(isAdmin)
	return msg
}

func GetPredictionMsg(chat_id int64, textMessage string, matchName string, matchID int, tag string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chat_id, textMessage)
	msg.ReplyMarkup = GetBetOnMatchKeyboard(matchName, matchID, tag)
	return msg
}

func GetConfirmFormMsg(chat_id int64, confirm_predictions model.ConfirmPrediction) tgbotapi.MessageConfig {
	match_id := strconv.Itoa(confirm_predictions.Match_id)
	msg := tgbotapi.NewMessage(chat_id, confirm_predictions.TextMsg)
	msg.ReplyMarkup = GetConfirmFormKeyboard(match_id, confirm_predictions.Bet, confirm_predictions.Tag)
	return msg
}

// Keyboard

func GetMainMenuKeyboard(userIsAdmin bool) tgbotapi.InlineKeyboardMarkup {

	var keyboard tgbotapi.InlineKeyboardMarkup

	if userIsAdmin {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Новый турнир", "create_tournament"),
				tgbotapi.NewInlineKeyboardButtonData("➕ Новые матчи", "create_matches"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎯 Добавить результат", "add_results"),
				tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить турнир", "finish_tournament"),
				tgbotapi.NewInlineKeyboardButtonData("*️⃣ Получить ИД матчей", "get_match_ids"),
			),
		)
	}

	keyboard.InlineKeyboard = append(
		keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 Мои прогнозы", "user_predictions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Сделать ставку", "active_matches"),
		),
	)

	return keyboard
}

func GetBetOnMatchKeyboard(matchName string, matchID int, tag string) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0️⃣", tag+"_prediction_"+matchName+"_"+strconv.Itoa(matchID)+"_0-2"),
			tgbotapi.NewInlineKeyboardButtonData("1️⃣", tag+"_prediction_"+matchName+"_"+strconv.Itoa(matchID)+"_1-2"),
			tgbotapi.NewInlineKeyboardButtonData("2️⃣", tag+"_prediction_"+matchName+"_"+strconv.Itoa(matchID)+"_2-0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎯 Win team 1", tag+"_prediction_"+matchName+"_"+strconv.Itoa(matchID)+"_1"),
			tgbotapi.NewInlineKeyboardButtonData("🎯 Win team 2", tag+"_prediction_"+matchName+"_"+strconv.Itoa(matchID)+"_2"),
		),
	)
}

func GetConfirmFormKeyboard(matchID string, bet string, tag string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да", tag+"_prediction_"+matchID+"_"+bet+"_y"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", tag+"_prediction_"+matchID+"_"+bet+"_n"),
		),
	)
}
