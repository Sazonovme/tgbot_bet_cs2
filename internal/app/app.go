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
	if update.Message != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		userData := PrepareUserData(update)

		switch {
		case userData.TextMsg == "/create-tournament":
			a.Handler.CreateTournament(ctx, userData)
		case userData.TextMsg == "/create-match":
			a.Handler.CreateMatch(ctx, userData)
		case userData.TextMsg == "/add-result":
			a.Handler.AddMatchResult(ctx, userData)
		case userData.TextMsg == "/finish-tournament":
			a.Handler.FinishTournament(ctx, userData)
		case userData.TextMsg == "/my-predictions":
			a.Handler.MyPredictions(ctx, userData)
		case strings.Contains(update.Message.Text, "/match"):
			a.Handler.MakePrediction(ctx, userData)
		case userData.TextMsg == "/start":
			a.Handler.Start(ctx, userData)
		case userData.TextMsg == "/stop":
			a.Handler.Stop(ctx, userData)
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
