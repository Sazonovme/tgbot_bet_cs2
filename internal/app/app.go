package app

import (
	"RushBananaBet/internal/handler"
	"RushBananaBet/internal/model"
	"RushBananaBet/pkg/logger"
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

	a.StartPolling()

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

	// === Inline-кнопки (CallbackQuery) ===
	if update.CallbackQuery != nil {
		callback := update.CallbackQuery
		userData := PrepareUserDataFromCallback(callback)

		switch {
		case strings.HasPrefix(callback.Data, "create_tournament"):
			a.Handler.CreateTournament(ctx, userData)
		case strings.HasPrefix(callback.Data, "create_match"):
			a.Handler.CreateMatch(ctx, userData)
		case strings.HasPrefix(callback.Data, "add_result"):
			a.Handler.AddMatchResult(ctx, userData)
		case strings.HasPrefix(callback.Data, "finish_tournament"):
			a.Handler.FinishTournament(ctx, userData)
		case strings.HasPrefix(callback.Data, "my_predictions"):
			a.Handler.MyPredictions(ctx, userData)
		case strings.HasPrefix(callback.Data, "match_"):
			a.Handler.MakePrediction(ctx, userData)
		case strings.HasPrefix(callback.Data, "bet_"):
			a.Handler.HandleBetSelection(ctx, userData, callback)
		case strings.HasPrefix(callback.Data, "cancel_"):
			a.Handler.HandleCancel(ctx, userData, callback)
		default:
			a.Handler.UnknownCallback(ctx, userData, callback)
		}
		return

	} else if update.Message != nil {

		userData := PrepareUserData(update)

		// === Обычные команды ===
		switch userData.TextMsg {
		case "/start":
			a.Handler.Start(ctx, userData)
		case "/stop":
			a.Handler.Stop(ctx, userData)
		default:
			a.Handler.UnknownCommand(ctx, userData)
		}
	}
}

func PrepareUserData(update tgbotapi.Update) *model.User {
	return &model.User{
		Chat_id:    update.Message.Chat.ID,
		Username:   update.Message.From.UserName,
		First_name: update.Message.From.FirstName,
		Last_name:  update.Message.From.LastName,
		TextMsg:    update.Message.Text,
	}
}

func PrepareUserDataFromCallback(callback *tgbotapi.CallbackQuery) *model.User {
	return &model.User{
		Chat_id:    callback.Message.Chat.ID,
		Username:   callback.Message.From.UserName,
		First_name: callback.Message.From.FirstName,
		Last_name:  callback.Message.From.LastName,
		TextMsg:    callback.Message.Text,
	}
}
