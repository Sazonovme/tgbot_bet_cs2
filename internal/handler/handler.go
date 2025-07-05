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
	GetEventFinishTable(ctx context.Context) ([]model.FinishTable, error)
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

func (h *Handler) CreateEvent() {

}

func (h *Handler) AddResult() {

}

func (h *Handler) FinishTournament() {

}

func (h *Handler) MyPredictions() {

}

func (h *Handler) MakePrediction() {

}
