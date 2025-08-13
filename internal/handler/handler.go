package handler

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	safemap "RushBananaBet/internal/safeMap"
	"RushBananaBet/internal/ui"
	"context"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var UserSessionsMap *safemap.SafeMap

type Handler struct {
	BotApi  *tgbotapi.BotAPI
	Service Service
}

type Service interface {
	CreateTournament(ctx context.Context, userData *model.User) error
	CreateMatch(ctx context.Context, userData *model.User) error
	AddMatchResult(ctx context.Context, userData *model.User) error
	GetTournamentFinishTable(ctx context.Context) (*[]model.TournamentFinishTable, *model.ScoreFinishTable, error)
	GetActiveMatches(ctx context.Context) (*[]model.Match, error)
	GetUserPredictions(ctx context.Context, username string) (*[]model.UserPrediction, error)
	AddUserPrediction(ctx context.Context, userData *model.User) error
	AddNewUser(ctx context.Context, user *model.User) (err error, isExist bool)
	DeactivateUser(ctx context.Context, chat_id int64) error
}

func init() {
	UserSessionsMap = safemap.NewSafeMap()
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
		_, err = sendMsg(h.BotApi, userData.Chat_id, "Привет, теперь ты участник закрытого клуба петушков", nil)
		if err != nil {
			h.Service.DeactivateUser(ctx, userData.Chat_id)
			return
		}
		keyboard:= ui.PaintMainMenu(model.IsAdmin(userData.Username))
		msg, err := sendMsg(h.BotApi, userData.Chat_id, "Главное меню:", keyboard)
		if err != nil {
			logger.Error("Err start()", "handler-Start()", err)
			return
		}
		//LastMsgType.Set(userData.Chat_id, msg.MessageID)
	}
}

func (h *Handler) Stop(ctx context.Context, userData *model.User) {
	sendMsg(h.BotApi, userData.Chat_id, "gg, ты больше не участник, так даже лучше, ТАКИЕ писькотрясы нам не нужны", nil)
	h.Service.DeactivateUser(ctx, userData.Chat_id)
}

func (h *Handler) CreateTournament(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", nil)
		return
	}
	h.Service.CreateTournament(ctx, userData)
}

func (h *Handler) CreateMatch(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", nil)
		return
	}
	h.Service.CreateMatch(ctx, userData)
}

func (h *Handler) AddMatchResult(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", nil)
		return
	}
	h.Service.AddMatchResult(ctx, userData)
}

func (h *Handler) FinishTournament(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", nil)
		return
	}
	h.Service.GetTournamentFinishTable(ctx)
}

// Все матчи + сделать прогноз

func (h *Handler) GetActiveMatches(ctx context.Context, userData *model.User) {

	// Очищаем чат от старых сообщений
	// Если сообщение отправлено больше чем 48ч назад удалить не можем
	// Либо удаляем сообщение либо клавиатуру
	// lastMessagesIDs, sendAt, _, ok := UserSessionsMap.Get(userData.Chat_id)
	// if !ok {
	// 	logger.Error("No value in map", "handler-GetActiveMatches()", nil)
	// 	return
	// }

	// if time.Since(sendAt) > 47*time.Hour + 55 * time.Minute {
	// 	for _, lastMsgID := range lastMessagesIDs {
	// 		deleteKeyboard(h.BotApi, userData.Chat_id, lastMsgID)
	// 	}	
	// } else {
	// 	for _, lastMsgID := range lastMessagesIDs {
	// 		deleteMsg(h.BotApi, userData.Chat_id, lastMsgID)
	// 	}		
	// }

	matches, err := h.Service.GetActiveMatches(ctx)
	if err != nil {
		logger.Error("Error get active matches in handler", "handler-GetActiveMatches()", err)
		return
	}

	msg, err := sendMsg(h.BotApi, userData.Chat_id, "🔽 ВСЕ АКТИВНЫЕ МАТЧИ 🔽")
	if err != nil {
		logger.Error("Error send msg", "handler-GetActiveMatches()", err)
		return
	}

	for _, match := range *matches {

		var MessageIDs []int

		keyboard:= ui.PaintButtonsForBetOnMatch(match.Id)
		msgText := match.Name + "\n" + "Выберите точный счет для команды " + match.Team1 + " или Win для ставки на победу команды"
		msg, err := sendMsg(h.BotApi, userData.Chat_id, msgText, keyboard)
		if err != nil {
			logger.Error("Err GetActiveMatches()", "handler-GetActiveMatches()", err)
			return
		}

		MessageIDs = append(MessageIDs, msg.MessageID)
	} 

	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, MessageIDs, "active_matches")
}

func (h *Handler) MakePrediction(ctx context.Context, userData *model.User) {
	h.Service.AddUserPrediction(ctx, userData)
}

func (h *Handler) MyPredictions(ctx context.Context, userData *model.User) {
	userPredictions, err := h.Service.GetUserPredictions(ctx, userData.Username)
	if err != nil {
		logger.Error("Dont recive user predictions", "handler-MyPredictions()", err)
		return
	}
}



func (h *Handler) UnknownCommand(ctx context.Context, userData *model.User) {
	sendMsg(h.BotApi, userData.Chat_id, "Неизвестная команда", nil)
}

func (h *Handler) HandleBackTo(ctx context.Context, userData *model.User, callback *tgbotapi.CallbackQuery) {
	
	pointMenu := strings.Replace(callback.Data, "back_to_", "") 
	switch pointMenu {
		case 
	}
	callback.Data
}
