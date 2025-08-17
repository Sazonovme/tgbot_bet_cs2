package ui

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PaintMainMenu(userIsAdmin bool) tgbotapi.InlineKeyboardMarkup {

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
			),
		)
	}

	keyboard.InlineKeyboard = append(
		keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 Мои прогнозы", "my_predictions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Сделать ставку", "active_matches"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Завершить участие", "get_out"),
		),
	)

	return keyboard
}

func PaintButtonsForBetOnMatch(matchName string, matchID int, tag string) tgbotapi.InlineKeyboardMarkup {

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

func PaintConfirmForm(tag string, matchID string, bet string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да", tag+"_prediction_"+matchID+"_"+bet+"_y"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", tag+"_prediction_"+matchID+"_"+bet+"_n"),
		),
	)
}

//.  🔙

func PaintUserPredictions() {

}

// func BuildKeyboard(username string) tgbotapi.ReplyKeyboardMarkup {
// 	var rows [][]tgbotapi.KeyboardButton

// 	// Админские кнопки
// 	if user.IsAdmin(username) {
// 		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("➕ Добавить ивент"),
// 		))
// 		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("🎯 Добавить результат"),
// 		))
// 		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
// 			tgbotapi.NewKeyboardButton("🏁 Завершить турнир"),
// 		))
// 	}

// 	// Пользовательская кнопка
// 	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("📄 Мои ставки"),
// 	))

// 	// Кнопки матчей
// 	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
// 		tgbotapi.NewKeyboardButton("⚔️ Матч 1"),
// 		tgbotapi.NewKeyboardButton("⚔️ Матч 2"),
// 	))

// 	return tgbotapi.NewReplyKeyboard(rows...)
// }
