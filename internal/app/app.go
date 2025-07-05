package app

import (
	"RushBananaBet/internal/handler"
	"RushBananaBet/internal/model"
	"RushBananaBet/pkg/logger"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	botApi  *tgbotapi.BotAPI
	Handler handler.Handler
}

func NewApp(botToken string, handler handler.Handler) App {
	botApi, err := tgbotapi.NewBotAPI(botToken)
	handler.BotApi = botApi
	if err != nil {
		logger.Fatal("Error creating newBot", "main-main()", err)
	}
	return App{
		botApi:  botApi,
		Handler: handler,
	}
}

func (a *App) Start(stop chan os.Signal) {
	logger.Info("Bot started", "app-(*Bot)Start()", nil)

	a.StartPolling()

	<-stop
	logger.Info("Bot stoped", "app-(*Bot)Start()", nil)

	a.botApi.StopReceivingUpdates()
}

func (a *App) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.botApi.GetUpdatesChan(u)

	for update := range updates {
		go a.RouteUpdate(update)
	}
}

func (a *App) RouteUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		data := PrepareData(update)
		switch {
		case update.Message.Text == "/start":
			a.Handler.Start(data)
		case update.Message.Text == "/create-event":
			// ds
		case update.Message.Text == "/add-result":
			//ds
		case update.Message.Text == "/finish-tournament":
			// ыв
		case update.Message.Text == "/my-predictions":
			//ds
		case strings.Contains(update.Message.Text, "/match"):
			// dsd
		}
	}
}

func PrepareData(update tgbotapi.Update) model.HandlerData {
	return model.HandlerData{
		ChatID:   update.Message.Chat.ID,
		UserName: update.Message.From.UserName,
		Text:     update.Message.Text,
	}
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
