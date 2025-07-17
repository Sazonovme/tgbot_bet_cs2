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
	CreateTournament(ctx context.Context, name_tournament string) error
	CreateMatch(ctx context.Context, match *model.Match) error
	AddMatchResult(ctx context.Context, result string, match_id int) error
	GetTournamentFinishTable(ctx context.Context) (*[]model.TournamentFinishTable, *model.ScoreFinishTable, error)
	GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, prediction *model.UserPrediction) error
	AddNewUser(ctx context.Context, user *model.User) error
	DeactivateUser(ctx context.Context, chat_id int64) error
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
		h.Service.AddNewUser(ctx, user)
	} else {
		h.Service.DeactivateUser(ctx, user.Chat_id)
	}
}

func (h *Handler) Stop(user *model.User, ctx context.Context) {
	msg := tgbotapi.NewMessage(user.Chat_id, "gg, ты больше не участник, так даже лучше, лузеры нам не нужны")
	h.BotApi.Send(msg)
	h.Service.DeactivateUser(ctx, user.Chat_id)
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
