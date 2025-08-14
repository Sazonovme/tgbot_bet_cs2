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
	err, _ := h.Service.AddNewUser(ctx, userData)
	if err != nil {
		logger.Error("Err start()", "handler-Start()", err)
		return
	}

	_, err = sendMsg(h.BotApi, userData.Chat_id, "Привет, теперь ты участник закрытого клуба петушков", tgbotapi.InlineKeyboardMarkup{})
	if err != nil {
		h.Service.DeactivateUser(ctx, userData.Chat_id)
		return
	}
	keyboard := ui.PaintMainMenu(model.IsAdmin(userData.Username))
	msg, err := sendMsg(h.BotApi, userData.Chat_id, "Главное меню:", keyboard)
	if err != nil {
		logger.Error("Err start()", "handler-Start()", err)
		return
	}
	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, []int{msg.MessageID}, "active_matches")

}

func (h *Handler) Stop(ctx context.Context, userData *model.User) {
	sendMsg(h.BotApi, userData.Chat_id, "gg, ты больше не участник, так даже лучше, ТАКИЕ писькотрясы нам не нужны", tgbotapi.InlineKeyboardMarkup{})
	h.Service.DeactivateUser(ctx, userData.Chat_id)
}

func (h *Handler) CreateTournament(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}
	h.Service.CreateTournament(ctx, userData)
}

func (h *Handler) CreateMatch(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}
	h.Service.CreateMatch(ctx, userData)
}

func (h *Handler) AddMatchResult(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}
	h.Service.AddMatchResult(ctx, userData)
}

func (h *Handler) FinishTournament(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
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

	_, err = sendMsg(h.BotApi, userData.Chat_id, "🔽 ВСЕ АКТИВНЫЕ МАТЧИ 🔽", tgbotapi.InlineKeyboardMarkup{})
	if err != nil {
		logger.Error("Error send msg", "handler-GetActiveMatches()", err)
		return
	}

	var MessageIDs []int

	for _, match := range *matches {

		keyboard := ui.PaintButtonsForBetOnMatch(match.Name, match.Id, "confirm")
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

func (h *Handler) ConfirmPrediction(ctx context.Context, userData *model.User) {

	// confirm_prediction_[matchName]_[matchID]_[bet]
	arr := strings.Split(userData.CallbackData, "_")

	// Расшифровка ставки
	betTxt := ""
	if arr[4] == "1" || arr[4] == "2" {
		betTxt = "Победа команды " + arr[4]
	} else {
		betTxt = "Точный счет " + arr[4]
	}

	// Команды
	teams := strings.Split(arr[2], "vs")

	// Итоговое сообщение
	textMessage := "Матч: " + teams[0] + " vs " + teams[1] + "\n" + "Ваша ставка: " + betTxt + "\n" + "Подтвердить ставку?"

	keyboard := ui.PaintConfirmForm(arr[2], arr[3])

	msg, err := sendMsg(h.BotApi, userData.Chat_id, textMessage, keyboard)
	if err != nil {
		logger.Error("Error send msg", "handler-ConfirmPrediction()", err)
		return
	}

	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, []int{msg.MessageID}, "confirm_form")
}

func (h *Handler) MakePrediction(ctx context.Context, userData *model.User) {

	// make_prediction_[matchID]_[bet]_[y/n]
	arr := strings.Split(userData.CallbackData, "_")
	if arr[4] == "n" {
		h.GetActiveMatches(ctx, userData)
		return
	}

	err := h.Service.AddUserPrediction(ctx, userData)
	if err != nil {
		logger.Error("Err add user prediction", "handler-MakePrediction()", err)
		return
	}

	// Переводим в главное меню
	_, err = sendMsg(h.BotApi, userData.Chat_id, "✅ Ставка успешно сделана ✅", tgbotapi.InlineKeyboardMarkup{})
	if err != nil {
		logger.Error("Err send msg", "handler-MakePrediction()", err)
		return
	}

	keyboard := ui.PaintMainMenu(model.IsAdmin(userData.Username))
	msg, err := sendMsg(h.BotApi, userData.Chat_id, "Главное меню:", keyboard)
	if err != nil {
		logger.Error("Err start()", "handler-Start()", err)
		return
	}

	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, []int{msg.MessageID}, "main_menu")
}

func (h *Handler) MyPredictions(ctx context.Context, userData *model.User) {

	// Необходимо построить таблицу матчей
	// 1. Матчи дата которых уже истекла должны быть без кнопок
	// 2. Матчи до даты с кнопками
	// 3. Матчи с результатом должны быть сразу с количеством баллов

	userPredictions, err := h.Service.GetUserPredictions(ctx, userData.Username)
	if err != nil {
		logger.Error("Dont recive user predictions", "handler-MyPredictions()", err)
		return
	}

	var MessageIDs []int

	for _, prediction := range *userPredictions {

		readablePrediction := ""
		if prediction.Prediction == "1" || prediction.Prediction == "2" {
			readablePrediction = "Team " + prediction.Prediction + " win"
		} else {
			readablePrediction = prediction.Prediction
		}

		// Завершенные матчи
		if prediction.Result != "" {
			points := CalcPointsForBet(prediction.Prediction, prediction.Result)
			txtMsg := "✅ Матч завершен ✅" + "\n" + prediction.Match_Name + "\n" + "Счет матча: " + prediction.Result + "\n" + "Твоя ставка: " + readablePrediction + "\n" + "Заработанные баллы: " + points
			msg, err := sendMsg(h.BotApi, userData.Chat_id, txtMsg, tgbotapi.InlineKeyboardMarkup{})
			if err != nil {
				logger.Error("Error send msg", "handler-MyPredictions()", err)
				return
			}
			MessageIDs = append(MessageIDs, msg.MessageID)
			continue
		}

		// Текущие матчи
		if prediction.DateMatch.Before(time.Now()) {
			txtMsg := "🔴 Текущий матч 🔴" + "\n" + prediction.Match_Name + "\n" + "Твоя ставка: " + readablePrediction
			msg, err := sendMsg(h.BotApi, userData.Chat_id, txtMsg, tgbotapi.InlineKeyboardMarkup{})
			if err != nil {
				logger.Error("Error send msg", "handler-MyPredictions()", err)
				return
			}
			MessageIDs = append(MessageIDs, msg.MessageID)
			continue
		}

		// Будущие матчи
		txtMsg := "🔵 Матч еще не начался 🔵" + "\n" + prediction.Match_Name + "\n" + "Твоя ставка: " + readablePrediction
		keyboard := ui.PaintButtonsForBetOnMatch(prediction.Match_Name, int(prediction.Match_id), "change")
		msg, err := sendMsg(h.BotApi, userData.Chat_id, txtMsg, keyboard)
		if err != nil {
			logger.Error("Error send msg", "handler-ConfirmPrediction()", err)
			return
		}
		MessageIDs = append(MessageIDs, msg.MessageID)
	}

	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, MessageIDs, "my_predictions")
}

func (h *Handler) UnknownCommand(ctx context.Context, userData *model.User) {
	//sendMsg(h.BotApi, userData.Chat_id, "Неизвестная команда", nil)
}

func (h *Handler) HandleBackTo(ctx context.Context, userData *model.User, callback *tgbotapi.CallbackQuery) {

	// pointMenu := strings.Replace(callback.Data, "back_to_", "")
	// switch pointMenu {
	// 	case
	// }
	// callback.Data
}
