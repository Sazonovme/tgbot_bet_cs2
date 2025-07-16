package handler

import (
	"RushBananaBet/internal/model"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	BotApi  *tgbotapi.BotAPI
	Service Service
}

type Service interface {
	CreateEvent(ctx context.Context, event model.Event) error
	AddResultToEvent(ctx context.Context, result string) error
	GetEventFinishTable(ctx context.Context) ([]model.EventFinishTable, model.ScoreFinishTable, error)
	GetUserPredictions(ctx context.Context, username string) ([]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) Start(data model.HandlerData) {
	msg := tgbotapi.NewMessage(data.ChatID, "Привет, это бот для ставок на КС2, lets go")
	h.BotApi.Send(msg)
}

func (h *Handler) CreateEvent(data model.HandlerData) {

}

func (h *Handler) AddResult(data model.HandlerData) {

}

func (h *Handler) FinishTournament(data model.HandlerData) {

}

func (h *Handler) MyPredictions(data model.HandlerData) {

}

func (h *Handler) MakePrediction(data model.HandlerData) {

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
