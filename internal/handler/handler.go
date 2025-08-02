package handler

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	"RushBananaBet/internal/ui"
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
	AddNewUser(ctx context.Context, user *model.User) (err error, isExist bool)
	DeactivateUser(ctx context.Context, chat_id int64) error
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) Start(ctx context.Context, userData *model.User) {
	err, isExist := h.Service.AddNewUser(ctx, userData)
	if err != nil {
		logger.Error("Err start()", "handler-Start()", err)
		return
	}
	if !isExist {
		err2 := sendMsg(h.BotApi, userData.Chat_id, "Привет, теперь ты участник закрытого клуба петушков")
		if err2 != nil {
			h.Service.DeactivateUser(ctx, userData.Chat_id)
			return
		}
		ui.PaintMainMenu(model.IsAdmin(userData.Username))
	}
}

func (h *Handler) Stop(ctx context.Context, userData *model.User) {
	sendMsg(h.BotApi, userData.Chat_id, "gg, ты больше не участник, так даже лучше, ТАКИЕ писькотрясы нам не нужны")
	h.Service.DeactivateUser(ctx, userData.Chat_id)
}

func (h *Handler) CreateTournament(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов")
		return
	}
	h.Service.CreateTournament(ctx, userData)
}

func (h *Handler) CreateMatch(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов")
		return
	}
	h.Service.CreateMatch(ctx, userData)
}

func (h *Handler) AddMatchResult(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов")
		return
	}
	h.Service.AddMatchResult(ctx, userData)
}

func (h *Handler) FinishTournament(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов")
		return
	}
	h.Service.GetTournamentFinishTable(ctx)
}

func (h *Handler) MyPredictions(ctx context.Context, userData *model.User) {
	_, err := h.Service.GetUserPredictions(ctx, userData.Username)
	if err != nil {
		logger.Error("Dont recive user predictions", "handler-MyPredictions()", err)
		return
	}
}

func (h *Handler) MakePrediction(ctx context.Context, userData *model.User) {
	h.Service.AddUserPrediction(ctx, userData)
}

func (h *Handler) UnknownCommand(ctx context.Context, userData *model.User) {
	sendMsg(h.BotApi, userData.Chat_id, "Неизвестная команда")
}
