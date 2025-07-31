package ui

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func PaintMainMenu(userIsAdmin bool) ([][]tgbotapi.KeyboardButton, string) {

	txtMessage := "Ты в главном меню, выбери один из пунктов:"

	var rows [][]tgbotapi.KeyboardButton

	if userIsAdmin {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Новый турнир"),
			tgbotapi.NewKeyboardButton("➕ Новый матч"),
		))
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🎯 Добавить результат"),
			tgbotapi.NewKeyboardButton("🏁 Завершить турнир"),
		))
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("📄 Мои прогнозы"),
	))
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⚔️ Сделать прогноз"),
	))
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Завершить участие"),
	))

	return rows, txtMessage
}

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
