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
				tgbotapi.NewInlineKeyboardButtonData("➕ Новый турнир", "new_tournament"),
				tgbotapi.NewInlineKeyboardButtonData("➕ Новый матч", "new_match"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎯 Добавить результат", "add_result"),
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

func PaintButtonsForBetOnMatch(matchID int) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0️⃣", "make_prediction_"+strconv.Itoa(matchID)+"_0"),
			tgbotapi.NewInlineKeyboardButtonData("1️⃣", "make_prediction_"+strconv.Itoa(matchID)+"_1"),
			tgbotapi.NewInlineKeyboardButtonData("2️⃣", "make_prediction_"+strconv.Itoa(matchID)+"_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎯 Win team 1", "make_prediction_"+strconv.Itoa(matchID)+"_1win"),
			tgbotapi.NewInlineKeyboardButtonData("🎯 Win team 2", "make_prediction_"+strconv.Itoa(matchID)+"_2win"),
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
