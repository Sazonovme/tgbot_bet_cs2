package handler

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	safemap "RushBananaBet/internal/safeMap"
	"RushBananaBet/internal/ui"
	"context"
	"errors"
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
	AddUserPrediction(ctx context.Context, userData *model.User, chat_id int64) error
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

func (h *Handler) Start(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

	err, _ := h.Service.AddNewUser(ctx, userData)
	if err != nil {
		logger.Error("Err start()", "handler-Start()", err)
		return
	}

	err = openMainMenu(h.BotApi, "Привет, теперь ты участник закрытого клуба петушков", userData)
	if err != nil {
		logger.Error("Err start()", "handler-Start()", err)
		return
	}

}

func (h *Handler) Stop(ctx context.Context, update *tgbotapi.Update) {

	sendMsg(h.BotApi, update.Message.Chat.ID, "gg, ты больше не участник, так даже лучше, ТАКИЕ писькотрясы нам не нужны", tgbotapi.InlineKeyboardMarkup{})
	err := h.Service.DeactivateUser(ctx, update.Message.Chat.ID)
	if err != nil {
		logger.Error("Err deactivate user", "handler-Stop()", err)
		return
	}
}

func (h *Handler) CreateTournamentMessage(ctx context.Context, update *tgbotapi.Update) {
	msg, _ := sendMsg(h.BotApi, update.CallbackQuery.Message.Chat.ID, "Отправьте название турнира", tgbotapi.InlineKeyboardMarkup{})
	UserSessionsMap.ChangeLastMessages(update.CallbackQuery.Message.Chat.ID, []int{msg.MessageID}, "create_tournament")
}

func (h *Handler) CreateTournament(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

	if !userData.IsAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}

	// Валидация
	if userData.TextMsg == "" {
		sendMsg(h.BotApi, userData.Chat_id, "Формат команды /create-tournament [название]", tgbotapi.InlineKeyboardMarkup{})
		return
	}

	err := h.Service.CreateTournament(ctx, userData)
	if err != nil {
		logger.Error("Err create tournament", "handler-CreateTournament()", err)
		// Возврат в главное меню
		err := openMainMenu(h.BotApi, err.Error(), userData)
		if err != nil {
			logger.Error("Open menu err", "handler-CreateMatches()", err)
			return
		}
		return
	}

	// Возвращаем в главное меню
	err = openMainMenu(h.BotApi, "✅ Туринр успешно создан ✅", userData)
	if err != nil {
		logger.Error("Err create tournament", "handler-CreateTournament()", err)
		return
	}
}

func (h *Handler) CreateMatchesMessage(ctx context.Context, update *tgbotapi.Update) {
	msg, _ := sendMsg(h.BotApi, update.CallbackQuery.Message.Chat.ID, "Отправьте матчи в формате: [t1vst2]_[2025-08-16 15:00]#...", tgbotapi.InlineKeyboardMarkup{})
	UserSessionsMap.ChangeLastMessages(update.CallbackQuery.Message.Chat.ID, []int{msg.MessageID}, "create_matches")
}

func (h *Handler) CreateMatch(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

	if !userData.IsAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}

	errText := validateMatches(userData.TextMsg)
	if errText != "" {
		logger.Error("Validation matches error", "handler-CreateMatches()", errors.New(errText))
		// Возврат в главное меню
		err := openMainMenu(h.BotApi, errText, userData)
		if err != nil {
			logger.Error("Open menu err", "handler-CreateMatches()", err)
			return
		}
		return
	}

	err := h.Service.CreateMatch(ctx, userData)
	if err != nil {
		logger.Error("Err create matches", "handler-CreateMatches()", err)
		// Возвращаем в главное меню
		err = openMainMenu(h.BotApi, err.Error(), userData)
		if err != nil {
			logger.Error("Open menu err", "handler-CreateMatches()", err)
			return
		}
		return
	}

	// Возвращаем в главное меню
	err = openMainMenu(h.BotApi, "✅ Матчи успешно созданы ✅", userData)
	if err != nil {
		logger.Error("Open menu err", "handler-CreateMatches()", err)
		return
	}
}

func (h *Handler) AddMatchesResultMessage(ctx context.Context, update *tgbotapi.Update) {
	msg, _ := sendMsg(h.BotApi, update.CallbackQuery.Message.Chat.ID, "Отправьте результаты в формате: [matchID]_[result]#...", tgbotapi.InlineKeyboardMarkup{})
	UserSessionsMap.ChangeLastMessages(update.CallbackQuery.Message.Chat.ID, []int{msg.MessageID}, "add_results")
}

func (h *Handler) AddMatchResult(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

	if !userData.IsAdmin {
		sendMsg(h.BotApi, userData.Chat_id, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}

	errText := validateMatchesResults(userData.TextMsg)
	if errText != "" {
		logger.Error("Validate matches reulst err", "handler-AddMatchResult()", errors.New(errText))
		// Возврат в главное меню
		err := openMainMenu(h.BotApi, errText, userData)
		if err != nil {
			logger.Error("Err open menu", "handler-AddMatchResult()", err)
			return
		}
		return
	}

	err := h.Service.AddMatchResult(ctx, userData)
	if err != nil {
		logger.Error("Add matches result err", "handler-AddMatchResult()", err)
		// Возврат в главное меню
		err := openMainMenu(h.BotApi, err.Error(), userData)
		if err != nil {
			logger.Error("Err open menu", "handler-AddMatchResult()", err)
			return
		}
		return
	}
}

func (h *Handler) FinishTournament(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

	if !userData.IsAdmin {
		sendMsg(h.BotApi, update.Message.Chat.ID, "Васылек, ниче не попутал? Иди гуляй, данная функция для администраторов", tgbotapi.InlineKeyboardMarkup{})
		return
	}
	h.Service.GetTournamentFinishTable(ctx)
}

func (h *Handler) GetActiveMatches(ctx context.Context, update *tgbotapi.Update) {

	matches, err := h.Service.GetActiveMatches(ctx)
	if err != nil {
		logger.Error("Error get active matches in handler", "handler-GetActiveMatches()", err)
		return
	}

	_, err = sendMsg(h.BotApi, update.CallbackQuery.Message.Chat.ID, "🔽 ВСЕ АКТИВНЫЕ МАТЧИ 🔽", tgbotapi.InlineKeyboardMarkup{})
	if err != nil {
		logger.Error("Error send msg", "handler-GetActiveMatches()", err)
		return
	}

	var MessageIDs []int

	for _, match := range *matches {

		keyboard := ui.PaintButtonsForBetOnMatch(match.Name, match.Id, "confirm")
		msgText := match.Name + "\n" + "Выберите точный счет для команды " + match.Team1 + " или Win для ставки на победу команды"
		msg, err := sendMsg(h.BotApi, update.CallbackQuery.Message.Chat.ID, msgText, keyboard)
		if err != nil {
			logger.Error("Err GetActiveMatches()", "handler-GetActiveMatches()", err)
			return
		}

		MessageIDs = append(MessageIDs, msg.MessageID)
	}

	UserSessionsMap.Delete(update.CallbackQuery.Message.Chat.ID)
	UserSessionsMap.Set(update.CallbackQuery.Message.Chat.ID, MessageIDs, "active_matches")
}

func (h *Handler) ConfirmPrediction(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

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

	tag := ""
	if arr[0] == "confirm" {
		tag = "make"
	} else {
		tag = "change"
	}

	keyboard := ui.PaintConfirmForm(tag, arr[3], arr[4])

	msg, err := sendMsg(h.BotApi, userData.Chat_id, textMessage, keyboard)
	if err != nil {
		logger.Error("Error send msg", "handler-ConfirmPrediction()", err)
		return
	}

	UserSessionsMap.Delete(userData.Chat_id)
	UserSessionsMap.Set(userData.Chat_id, []int{msg.MessageID}, "confirm_form")
}

func (h *Handler) MakePrediction(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

	// make_prediction_[matchID]_[bet]_[y/n]
	// change_prediction_[matchID]_[bet]_[y/n]
	arr := strings.Split(userData.CallbackData, "_")
	if arr[4] == "n" && arr[0] == "make" {
		h.GetActiveMatches(ctx, update)
		return
	} else if arr[4] == "n" && arr[0] == "change" {
		h.MyPredictions(ctx, update)
		return
	}

	err := h.Service.AddUserPrediction(ctx, userData, userData.Chat_id)
	if err != nil {
		logger.Error("Err add user prediction", "handler-MakePrediction()", err)
		return
	}

	// Переводим в главное меню
	txtMsg := ""
	if arr[0] == "make" {
		txtMsg = "✅ Ставка успешно сделана ✅"
	} else if arr[0] == "change" {
		txtMsg = "✅ Ставка успешно изменена ✅"
	}

	// Возврат в главное меню
	err = openMainMenu(h.BotApi, txtMsg, userData)
	if err != nil {
		logger.Error("Err open menu", "handler-AddMatchResult()", err)
		return
	}
}

func (h *Handler) MyPredictions(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

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

func (h *Handler) Help(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

	if !userData.IsAdmin {
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

func (h *Handler) GetMatchesIDs(ctx context.Context, update *tgbotapi.Update) {

	// Prepare data
	userData := PrepareUserData(update)

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

func (h *Handler) UnknownCommand(ctx context.Context, update *tgbotapi.Update) {
	//sendMsg(h.BotApi, userData.Chat_id, "Неизвестная команда", nil)
}

// func (h *Handler) HandleBackTo(ctx context.Context, userData *model.User, callback *tgbotapi.CallbackQuery) {

// 	// pointMenu := strings.Replace(callback.Data, "back_to_", "")
// 	// switch pointMenu {
// 	// 	case
// 	// }
// 	// callback.Data
// }
