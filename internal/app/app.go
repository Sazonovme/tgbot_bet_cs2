package app

import (
	"RushBananaBet/internal/handler"
	"RushBananaBet/internal/logger"
	"context"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	botApi  *tgbotapi.BotAPI
	Handler handler.Handler
}

func NewApp(botToken string, handler handler.Handler) App {
	botApi, err := tgbotapi.NewBotAPI(botToken)
	handler.BotApi = botApi
	if err != nil {
		logger.Fatal("Error creating newBot", "main-main()", err)
	}
	return App{
		botApi:  botApi,
		Handler: handler,
	}
}

func (a *App) Start(stop chan os.Signal) {
	logger.Info("Bot started", "app-(*Bot)Start()", nil)

	go a.StartPolling()

	<-stop
	logger.Info("Bot stoped", "app-(*Bot)Start()", nil)

	a.botApi.StopReceivingUpdates()
}

func (a *App) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.botApi.GetUpdatesChan(u)

	for update := range updates {
		go a.RouteUpdate(update)
	}
}

func (a *App) RouteUpdate(update tgbotapi.Update) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	updatePointer := &update

	// === Inline-кнопки (CallbackQuery) ===
	if update.CallbackQuery != nil {
		// callback := update.CallbackQuery
		// userData := PrepareUserDataFromCallback(callback)

		switch {
		case strings.HasPrefix(update.CallbackQuery.Data, "create_tournament"):
			a.Handler.CreateTournamentMessage(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "create_matches"):
			a.Handler.CreateMatchesMessage(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "add_results"):
			a.Handler.AddMatchesResultMessage(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "finish_tournament"):
			a.Handler.FinishTournament(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "active_matches"):
			a.Handler.GetActiveMatches(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "confirm_prediction"):
			a.Handler.ConfirmPrediction(ctx, updatePointer)
		// ДОБАВИТЬ
		// case strings.HasPrefix(update.CallbackQuery.Data, "change_prediction"):
		// 	a.Handler.ConfirmPrediction(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "my_predictions"):
			a.Handler.MyPredictions(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "make_prediction"):
			a.Handler.MakePrediction(ctx, updatePointer)
		// case strings.HasPrefix(callback.Data, "bet_"):
		// 	a.Handler.HandleBet(ctx, userData, callback)
		// case strings.HasPrefix(update.CallbackQuery.Data, "back_to_"):
		// 	a.Handler.HandleBackTo(ctx, userData, callback)
		default:
			a.Handler.UnknownCommand(ctx, updatePointer)
		}
		return

	} else if update.Message != nil {

		//userData := PrepareUserData(update)

		// === Обычные команды ===
		switch update.Message.Text {
		case "/start":
			a.Handler.Start(ctx, updatePointer)
		case "/stop":
			a.Handler.Stop(ctx, updatePointer)
		case "/help":
			a.Handler.Help(ctx, updatePointer)
		default:
			a.RouteText(ctx, updatePointer)
		}
	}
}

func (a *App) RouteText(ctx context.Context, update *tgbotapi.Update) {
	_, _, state, ok := handler.UserSessionsMap.Get(update.Message.Chat.ID)
	if !ok {
		a.Handler.UnknownCommand(ctx, update)
	}
	switch state {
	case "create_tournament":
		a.Handler.CreateTournament(ctx, update)
	case "create_matches":
		a.Handler.CreateMatch(ctx, update)
	case "add_results":
		a.Handler.AddMatchResult(ctx, update)
	}
}
