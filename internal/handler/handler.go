package handler

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	safemap "RushBananaBet/internal/safeMap"
	"RushBananaBet/internal/ui"
	"context"
	"strconv"
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
	GetMatchesIDs(ctx context.Context) (*[]model.Match, error)
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
	err := h.Service.DeactivateUser(ctx, userData.Chat_id)
	if err != nil {
		logger.Error("Err create tournament", "handler-Stop()", err)
		return
	}
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

func (h *Handler) GetActiveMatches(ctx context.Context, userData *model.User) {

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
	// change_prediction_[matchName]_[matchID]_[bet]
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
	textMessage := ""
	if arr[0] == "confirm" {
		textMessage = "Матч: " + teams[0] + " vs " + teams[1] + "\n" + "Ваша ставка: " + betTxt + "\n" + "Подтвердить ставку?"
	} else {
		textMessage = "Изменение ставки\n" + "Матч: " + teams[0] + " vs " + teams[1] + "\n" + "Новая ставка: " + betTxt + "\n" + "Подтвердить?"
	}

	keyboard := ui.PaintConfirmForm(arr[2], arr[3], arr[0])

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
	// change_prediction_[matchID]_[bet]_[y/n]
	arr := strings.Split(userData.CallbackData, "_")
	if arr[4] == "n" && arr[0] == "make" {
		h.GetActiveMatches(ctx, userData)
		return
	} else if arr[4] == "n" && arr[0] == "change" {
		h.MyPredictions(ctx, userData)
		return
	}

	err := h.Service.AddUserPrediction(ctx, userData)
	if err != nil {
		logger.Error("Err add user prediction", "handler-MakePrediction()", err)
		return
	}

	// Переводим в главное меню
	if arr[0] == "make" {
		_, err = sendMsg(h.BotApi, userData.Chat_id, "✅ Ставка успешно сделана ✅", tgbotapi.InlineKeyboardMarkup{})
	} else if arr[0] == "change" {
		_, err = sendMsg(h.BotApi, userData.Chat_id, "✅ Ставка успешно изменена ✅", tgbotapi.InlineKeyboardMarkup{})
	}
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

func (h *Handler) Help(ctx context.Context, userData *model.User) {
	isAdmin := model.IsAdmin(userData.Username)
	if !isAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}

	txtMsg :=
		`/create_tournament nameTournament
	/create_match name_date#name_date#
	/add_result matchid_result#matchid_result
	/finish_tournament`
	sendMsg(h.BotApi, userData.Chat_id, txtMsg, tgbotapi.InlineKeyboardMarkup{})
}

func (h *Handler) GetMatchesIDs(ctx context.Context, userData *model.User) {
	matches, err := h.Service.GetMatchesIDs(ctx)
	if err != nil {
		logger.Error("Err get match IDs in handler", "handler-GetMatchesIDs()", err)
		return
	}
	txtMsg := ""
	for _, match := range *matches {
		txtMsg += match.Date.Format("2006-01-02 15:04") + " " + match.Name + " " + strconv.Itoa(match.Id) + "\n"
	}

	sendMsg(h.BotApi, userData.Chat_id, txtMsg, tgbotapi.InlineKeyboardMarkup{})

	keyboard := ui.PaintMainMenu(model.IsAdmin(userData.Username))
	msg, err := sendMsg(h.BotApi, userData.Chat_id, "Главное меню:", keyboard)
	if err != nil {
		logger.Error("Err start()", "handler-GetMatchesIDs()", err)
		return
	}

	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, []int{msg.MessageID}, "main_menu")
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
