package handler

import (
	"RushBananaBet/internal/logger"
	"RushBananaBet/internal/model"
	"RushBananaBet/internal/ui"
	usersessions "RushBananaBet/internal/userSessions"
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var UserSessionsMap *usersessions.UserSessionMap

var ErrForUser = "Упс... Возникла ошибка - обратитесь к администратору"

type Handler struct {
	BotApi  *tgbotapi.BotAPI
	Service Service
}

type Service interface {

	//GENERAL
	AddNewUser(ctx context.Context, chat_id int64, user_id int64, username string) (isExist bool, err error)

	// ADMIN
	CreateTournament(ctx context.Context, name_tournament string) (added bool, err error)
	CreateMatches(ctx context.Context, matches []model.Match) error
	AddMatchResults(ctx context.Context, results []model.Result) error
	GetTournamentFinishTable(ctx context.Context) ([]model.TournamentFinishTable, model.ScoreFinishTable, error)
	GetActiveMatchesID(ctx context.Context) ([]model.Match, error)

	// USER
	GetActiveMatches(ctx context.Context) ([]model.Match, error)
	GetUserPredictions(ctx context.Context, chat_id int64) ([]model.UserPrediction, error)
	AddUpdateUserPrediction(ctx context.Context, chat_id int64, match_id int, prediction string) (inserted bool, err error)
}

func init() {
	UserSessionsMap = usersessions.NewUserSessionMap()
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) Start(ctx context.Context, update *tgbotapi.Update) {

	chat_id := update.Message.Chat.ID
	user_id := update.Message.From.ID
	username := update.Message.From.UserName

	isExist, err := h.Service.AddNewUser(ctx, chat_id, user_id, username)
	if err != nil {
		logger.Error("Err add new user", "handler-Start()", err)
		return
	}

	textMsg := ""
	if isExist {
		textMsg = "Ты и так участник закрытого клуба петушков"
	} else {
		textMsg = "Рады приветствовать тебя в закрытом клубе петушков"
	}

	userIsAdmin := model.IsAdmin(chat_id)
	message := ui.GetMainMenuMsg(chat_id, textMsg, userIsAdmin)

	sMsg, err := h.BotApi.Send(message)
	if err != nil {
		logger.Error("Err send msg in Start()", "handler-Start()", err)
		return
	}

	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
}

func (h *Handler) CreateTournamentMessage(ctx context.Context, update *tgbotapi.Update) {

	chat_id := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chat_id, "Отправьте название турнира")
	sMsg, err := h.BotApi.Send(msg)
	if err != nil {
		logger.Error("Err send create tournament msg", "handler-CreateTournamentMessage()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "create_tournament_msg")
}

func (h *Handler) CreateTournament(ctx context.Context, update *tgbotapi.Update) {

	textMsg := update.Message.Text
	chat_id := update.Message.Chat.ID

	name_tournament, err := prepareCreateTournamentData(textMsg)
	if err != nil {
		msg := tgbotapi.NewMessage(chat_id, "Ошибка. "+err.Error()+"\n Отправьте название еще раз")
		sMsg, err := h.BotApi.Send(msg)
		if err != nil {
			logger.Error("Err send create tournament msg", "handler-CreateTournament()", err)
			return
		}
		UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "create_tournament_msg")
		return
	}

	added, err := h.Service.CreateTournament(ctx, name_tournament)

	returnMessage := ""
	if err != nil {
		logger.Error("Err create tournament", "handler-CreateTournament()", err)
		returnMessage = "Ошибка " + err.Error()
	} else if added {
		returnMessage = "✅ Туринр успешно создан ✅"
	} else {
		returnMessage = "Уже есть активный турнир, сначала необходимо завершить прошлый турнир"
	}

	message := ui.GetMainMenuMsg(chat_id, returnMessage, model.IsAdmin(chat_id))
	sMsg, err := h.BotApi.Send(message)
	if err != nil {
		logger.Error("Err send create tournament msg", "handler-CreateTournament()", err)
		return
	}

	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
}

func (h *Handler) CreateMatchesMessage(ctx context.Context, update *tgbotapi.Update) {

	chat_id := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chat_id, "Отправьте матчи в формате: [t1_t2]-[16.08.2025 15:00]#...")
	sMsg, err := h.BotApi.Send(msg)
	if err != nil {
		logger.Error("Err send create tournament msg", "handler-CreateTournamentMessage()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "create_matches_msg")
}

func (h *Handler) CreateMatches(ctx context.Context, update *tgbotapi.Update) {

	textMsg := update.Message.Text
	chat_id := update.Message.Chat.ID

	matches, err := prepareCreateMatchesData(textMsg)
	if err != nil {
		msg := tgbotapi.NewMessage(chat_id, "Ошибка. "+err.Error()+"\n Отправьте еще раз")
		sMsg, err := h.BotApi.Send(msg)
		if err != nil {
			logger.Error("Err send create matches msg", "handler-CreateMatches()", err)
			return
		}
		UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "create_matches_msg")
		return
	}

	err = h.Service.CreateMatches(ctx, matches)

	returnMessage := ""
	if err != nil {
		returnMessage = "Ошибка: " + err.Error()
		logger.Error("Err create matches", "handler-CreateMatches()", err)
	} else {
		returnMessage = "✅ Матчи успешно добавлены ✅"
	}

	message := ui.GetMainMenuMsg(chat_id, returnMessage, model.IsAdmin(chat_id))
	sMsg, err := h.BotApi.Send(message)
	if err != nil {
		logger.Error("Err send msg in create matches", "handler-CreateMatches()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
}

func (h *Handler) AddMatchesResultMessage(ctx context.Context, update *tgbotapi.Update) {

	chat_id := update.CallbackQuery.Message.Chat.ID
	msg := tgbotapi.NewMessage(chat_id, "Отправьте результаты в формате [matchID]_[result]#...")
	sMsg, err := h.BotApi.Send(msg)
	if err != nil {
		logger.Error("Err send add results msg", "handler-AddMatchesResultMessage()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "add_results_msg")
}

func (h *Handler) GetActiveMatchesID(ctx context.Context, update *tgbotapi.Update) {

	chat_id := update.CallbackQuery.Message.Chat.ID
	textMsg := ""

	matches, err := h.Service.GetActiveMatchesID(ctx)
	if err != nil {
		logger.Error("Err get match IDs in handler", "handler-GetActiveMatchesID()", err)
		textMsg = "Err: " + err.Error()
	} else if len(matches) < 1 {
		textMsg = "❌ Матчей без результата не найдено ❌"
	} else {
		for _, match := range matches {
			textMsg += match.Date.Format("02.01.2006 15:04") + " " + match.Name + " " + strconv.Itoa(match.Id) + "\n"
		}
	}

	message := ui.GetMainMenuMsg(chat_id, textMsg, model.IsAdmin(chat_id))
	sMsg, err := h.BotApi.Send(message)
	if err != nil {
		logger.Error("Err send msg in GetActiveMatchesID()", "handler-GetActiveMatchesID()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
}

func (h *Handler) AddMatchResults(ctx context.Context, update *tgbotapi.Update) {

	textMsg := update.Message.Text
	chat_id := update.Message.Chat.ID

	results, err := prepareAddMatchResultsData(textMsg)
	if err != nil {
		msg := tgbotapi.NewMessage(chat_id, "Ошибка. "+err.Error()+"\n Отправьте еще раз")
		sMsg, err := h.BotApi.Send(msg)
		if err != nil {
			logger.Error("Err send add match results msg", "handler-AddMatchResult()", err)
			return
		}
		UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "add_results_msg")
		return
	}

	err = h.Service.AddMatchResults(ctx, results)

	returnMessage := ""
	if err != nil {
		returnMessage = "Ошибка: " + err.Error()
		logger.Error("Err add match results", "handler-AddMatchResults()", err)
	} else {
		returnMessage = "✅ Результаты матчей успешно добавлены ✅"
	}

	message := ui.GetMainMenuMsg(chat_id, returnMessage, model.IsAdmin(chat_id))
	sMsg, err := h.BotApi.Send(message)
	if err != nil {
		logger.Error("Err send msg in add match results", "handler-AddMatchResults()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
}

// TO DO
func (h *Handler) FinishTournament(ctx context.Context, update *tgbotapi.Update) {
	h.Service.GetTournamentFinishTable(ctx)
}

func (h *Handler) GetActiveMatches(ctx context.Context, update *tgbotapi.Update) {

	chat_id := update.CallbackQuery.Message.Chat.ID

	matches, err := h.Service.GetActiveMatches(ctx)

	textMsg := ""
	if err != nil {
		textMsg = ErrForUser
		logger.Error("Error get active matches in handler", "handler-GetActiveMatches()", err)
	} else if len(matches) < 1 {
		textMsg = "❌ Активных матчей не найдено ❌"
	}

	if textMsg != "" {
		message := ui.GetMainMenuMsg(chat_id, textMsg, model.IsAdmin(chat_id))
		sMsg, err := h.BotApi.Send(message)
		if err != nil {
			logger.Error("Err send msg in get active matches", "handler-GetActiveMatches()", err)
			return
		}
		UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
		return
	}

	var MessageIDs []int

	for _, match := range matches {

		msgText := match.Name + "\n" + "Выберите точный счет для команды " + match.Team1 + " или Win для ставки на победу команды"
		msg := ui.GetPredictionMsg(chat_id, msgText, match.Name, match.Id, "confirm")
		sMsg, err := h.BotApi.Send(msg)
		if err != nil {
			logger.Error("Err send msg in get active matches", "handler-GetActiveMatches()", err)
			return
		}
		MessageIDs = append(MessageIDs, sMsg.MessageID)
	}
	UserSessionsMap.Set(chat_id, MessageIDs, "active_matches")
}

func (h *Handler) GetUserPredictions(ctx context.Context, update *tgbotapi.Update) {

	chat_id := update.CallbackQuery.Message.Chat.ID

	userPredictions, err := h.Service.GetUserPredictions(ctx, chat_id)

	textMsg := ""
	if err != nil {
		textMsg = ErrForUser
		logger.Error("Err h.Service.GetUserPredictions()", "handler-GetUserPredictions()", err)
	} else if len(userPredictions) < 1 {
		textMsg = "❌ Активных ставок не найдено ❌"
	}

	if textMsg != "" {
		message := ui.GetMainMenuMsg(chat_id, textMsg, model.IsAdmin(chat_id))
		sMsg, err := h.BotApi.Send(message)
		if err != nil {
			logger.Error("Err send msg in get user predictions", "handler-GetUserPredictions()", err)
			return
		}
		UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
		return
	}

	var MessageIDs []int

	for _, prediction := range userPredictions {

		// Future matches
		if prediction.Result.String == "" {
			txtMsg := "🔵 Матч : " + prediction.Match_name + "\n" + "Твоя ставка: " + getReadableBet(prediction.Prediction)
			msg := ui.GetPredictionMsg(chat_id, txtMsg, prediction.Match_name, prediction.Match_id, "change")
			sMsg, err := h.BotApi.Send(msg)
			if err != nil {
				logger.Error("Err send msg in get user predictions", "handler-GetUserPredictions()", err)
				return
			}
			MessageIDs = append(MessageIDs, sMsg.MessageID)
		}

	}
	UserSessionsMap.Set(chat_id, MessageIDs, "user_predictions")
}

func (h *Handler) GetConfirmPrediction(ctx context.Context, update *tgbotapi.Update) {

	// 2 types data
	// confirm_prediction_[matchName]_[matchID]_[bet]
	// change_prediction_[matchName]_[matchID]_[bet]

	textMsg := update.CallbackQuery.Data
	chat_id := update.CallbackQuery.Message.Chat.ID

	confirmPrediction, err := prepareConfirmPredictionData(textMsg)
	if err != nil {
		message := ui.GetMainMenuMsg(chat_id, ErrForUser, model.IsAdmin(chat_id))
		sMsg, err := h.BotApi.Send(message)
		if err != nil {
			logger.Error("Err prepare confirm prediction data", "handler-ConfirmPrediction()", err)
			return
		}
		UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
		return
	}

	msg := ui.GetConfirmFormMsg(chat_id, confirmPrediction)
	sMsg, err := h.BotApi.Send(msg)
	if err != nil {
		logger.Error("Err send msg in confirm prediction", "handler-ConfirmPrediction()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "confirm_prediction")
}

func (h *Handler) ProcessingConfirmPrediction(ctx context.Context, update *tgbotapi.Update) {

	// 2 types data
	// Endconfirm_prediction_[matchID]_[bet]_[y/n]
	// Endchange_prediction_[matchID]_[bet]_[y/n]

	data := update.CallbackQuery.Data
	chat_id := update.CallbackQuery.Message.Chat.ID

	errFlag := false
	textMsg := ""

	preparedData, err := prepareProcessingConfirmPredictionData(data)
	if err != nil {
		logger.Error("Err prepareProcessingConfirmPredictionData()", "handler-ProcessingConfirmPrediction()", err)
		errFlag = true
		textMsg = ErrForUser
	}

	if !errFlag {

		// Rejection
		if !preparedData.Confirmed && preparedData.Tag == "Endconfirm" {
			h.GetActiveMatches(ctx, update)
			return
		} else if !preparedData.Confirmed && preparedData.Tag == "Endchange" {
			h.GetUserPredictions(ctx, update)
			return
		}

		added, err := h.Service.AddUpdateUserPrediction(ctx, chat_id, preparedData.Match_id, preparedData.Bet)
		if err != nil {
			logger.Error("Err h.Service.AddUpdateUserPrediction()", "handler-ProcessingConfirmPrediction()", err)
			textMsg = ErrForUser
		} else if added {
			textMsg = "✅ Ставка успешно совершена ✅"
		} else {
			textMsg = "✅ Ставка успешно изменена ✅"
		}
	}

	// Main menu
	message := ui.GetMainMenuMsg(chat_id, textMsg, model.IsAdmin(chat_id))
	sMsg, err := h.BotApi.Send(message)
	if err != nil {
		logger.Error("Err send msg in ProcessingConfirmPrediction()", "handler-ProcessingConfirmPrediction()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
}

// TO DO
func (h *Handler) GetUserPredictionsStats(ctx context.Context, update *tgbotapi.Update) {
	// userPredictions, err := h.Service.GetUserPredictions(ctx, userData.Username)
	// if err != nil {
	// 	logger.Error("Dont recive user predictions", "handler-MyPredictions()", err)
	// 	return
	// }

	// var MessageIDs []int

	// for _, prediction := range *userPredictions {

	// 	readablePrediction := ""
	// 	if prediction.Prediction == "1" || prediction.Prediction == "2" {
	// 		readablePrediction = "Team " + prediction.Prediction + " win"
	// 	} else {
	// 		readablePrediction = prediction.Prediction
	// 	}

	// 	// Завершенные матчи
	// 	if prediction.Result != "" {
	// 		points := CalcPointsForBet(prediction.Prediction, prediction.Result)
	// 		txtMsg := "✅ Матч завершен ✅" + "\n" + prediction.Match_Name + "\n" + "Счет матча: " + prediction.Result + "\n" + "Твоя ставка: " + readablePrediction + "\n" + "Заработанные баллы: " + points
	// 		msg, err := sendMsg(h.BotApi, userData.Chat_id, txtMsg, tgbotapi.InlineKeyboardMarkup{})
	// 		if err != nil {
	// 			logger.Error("Error send msg", "handler-MyPredictions()", err)
	// 			return
	// 		}
	// 		MessageIDs = append(MessageIDs, msg.MessageID)
	// 		continue
	// 	}

	// 	// Текущие матчи
	// 	if prediction.DateMatch.Before(time.Now()) {
	// 		txtMsg := "🔴 Текущий матч 🔴" + "\n" + prediction.Match_Name + "\n" + "Твоя ставка: " + readablePrediction
	// 		msg, err := sendMsg(h.BotApi, userData.Chat_id, txtMsg, tgbotapi.InlineKeyboardMarkup{})
	// 		if err != nil {
	// 			logger.Error("Error send msg", "handler-MyPredictions()", err)
	// 			return
	// 		}
	// 		MessageIDs = append(MessageIDs, msg.MessageID)
	// 		continue
	// 	}

	// 	// Будущие матчи
	// 	txtMsg := "🔵 Матч еще не начался 🔵" + "\n" + prediction.Match_Name + "\n" + "Твоя ставка: " + readablePrediction
	// 	keyboard := ui.PaintButtonsForBetOnMatch(prediction.Match_Name, int(prediction.Match_id), "change")
	// 	msg, err := sendMsg(h.BotApi, userData.Chat_id, txtMsg, keyboard)
	// 	if err != nil {
	// 		logger.Error("Error send msg", "handler-ConfirmPrediction()", err)
	// 		return
	// 	}
	// 	MessageIDs = append(MessageIDs, msg.MessageID)
	// }
}

func (h *Handler) UnknownCommand(ctx context.Context, chat_id int64, textMsg string) {
	message := ui.GetMainMenuMsg(chat_id, textMsg, model.IsAdmin(chat_id))
	sMsg, err := h.BotApi.Send(message)
	if err != nil {
		logger.Error("Err send msg in UnknownCommand()", "handler-UnknownCommand()", err)
		return
	}
	UserSessionsMap.Set(chat_id, []int{sMsg.MessageID}, "main_menu")
}
