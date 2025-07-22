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
	CreateTournament(ctx context.Context, userData *model.User) error
	CreateMatch(ctx context.Context, userData *model.User) error
	AddMatchResult(ctx context.Context, userData *model.User) error
	GetTournamentFinishTable(ctx context.Context) (*[]model.TournamentFinishTable, *model.ScoreFinishTable, error)
	GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, userData *model.User) error
	AddNewUser(ctx context.Context, user *model.User) error
	DeactivateUser(ctx context.Context, chat_id int64) error
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateTournament(ctx context.Context, userData *model.User) {
	err := middlewareAuth(h.BotApi, userData)
	if err == nil {
		h.Service.CreateTournament(ctx, userData)
	}
}

func (h *Handler) CreateMatch(ctx context.Context, userData *model.User) {
	err := middlewareAuth(h.BotApi, userData)
	if err == nil {
		h.Service.CreateMatch(ctx, userData)
	}
}

func (h *Handler) AddMatchResult(ctx context.Context, userData *model.User) {
	err := middlewareAuth(h.BotApi, userData)
	if err == nil {
		h.Service.AddMatchResult(ctx, userData)
	}
}

func (h *Handler) FinishTournament(ctx context.Context, userData *model.User) {
	err := middlewareAuth(h.BotApi, userData)
	if err == nil {
		h.Service.GetTournamentFinishTable(ctx)
	}
}

func (h *Handler) MyPredictions(ctx context.Context, userData *model.User) {
	h.Service.GetUserPredictions(ctx, userData.Username)
}

func (h *Handler) MakePrediction(ctx context.Context, userData *model.User) {
	h.Service.AddUserPrediction(ctx, userData)
}

func (h *Handler) Start(ctx context.Context, userData *model.User) {
	err := h.Service.AddNewUser(ctx, userData)
	err2 := sendMsg(h.BotApi, userData.Chat_id, "Привет, теперь ты участник закрытого клуба петушков")
	if err != nil || err2 != nil {
		h.Service.DeactivateUser(ctx, userData.Chat_id)
	}
}

func (h *Handler) Stop(ctx context.Context, userData *model.User) {
	sendMsg(h.BotApi, userData.Chat_id, "gg, ты больше не участник, так даже лучше, ТАКИЕ лохопеды нам не нужны")
	h.Service.DeactivateUser(ctx, userData.Chat_id)
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
