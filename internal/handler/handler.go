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
	Start(ctx context.Context, user *model.User) error
	Stop(ctx context.Context, chat_id int64) error
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

func (h *Handler) Start(user *model.User, ctx context.Context) {
	msg := tgbotapi.NewMessage(user.Chat_id, "Привет, это бот для ставок на КС2, lets go")
	_, err := h.BotApi.Send(msg)
	if err == nil {
		h.Service.Start(ctx, user)
	} else {
		h.Service.Stop(ctx, user.Chat_id)
	}
}

func (h *Handler) Stop(user *model.User, ctx context.Context) {
	msg := tgbotapi.NewMessage(user.Chat_id, "gg, ты больше не участник, так даже лучше, лузеры нам не нужны")
	h.BotApi.Send(msg)
	h.Service.Stop(ctx, user.Chat_id)
}

func (h *Handler) CreateEvent(user *model.User, ctx context.Context) {

}

func (h *Handler) AddResult(user *model.User, ctx context.Context) {

}

func (h *Handler) FinishTournament(user *model.User, ctx context.Context) {

}

func (h *Handler) MyPredictions(user *model.User, ctx context.Context) {

}

func (h *Handler) MakePrediction(user *model.User, ctx context.Context) {

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
