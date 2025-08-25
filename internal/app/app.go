package app

import (
	"RushBananaBet/internal/handler"
	"RushBananaBet/internal/logger"
	"context"
	"os"
	"strconv"
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

	// === Inline (CallbackQuery) ===
	if update.CallbackQuery != nil {
		switch {
		case strings.HasPrefix(update.CallbackQuery.Data, "create_tournament"):
			a.Handler.CreateTournamentMessage(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "create_matches"):
			a.Handler.CreateMatchesMessage(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "get_match_ids"):
			a.Handler.GetActiveMatchesID(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "add_results"):
			a.Handler.AddMatchesResultMessage(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "finish_tournament"):
			a.Handler.FinishTournament(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "active_matches"):
			a.Handler.GetActiveMatches(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "user_predictions"):
			a.Handler.GetUserPredictions(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "confirm") ||
			strings.HasPrefix(update.CallbackQuery.Data, "change"):
			a.Handler.GetConfirmPrediction(ctx, updatePointer)
		case strings.HasPrefix(update.CallbackQuery.Data, "Endconfirm") ||
			strings.HasPrefix(update.CallbackQuery.Data, "Endchange"):
			a.Handler.ProcessingConfirmPrediction(ctx, updatePointer)
		default:
			a.Handler.UnknownCommand(ctx, update.CallbackQuery.Message.Chat.ID, "🔴 Invalid command 🔴")
		}
		return

	} else if update.Message != nil {

		// === Regular commands ===
		switch update.Message.Text {
		case "/start":
			a.Handler.Start(ctx, updatePointer)
		default:
			a.RouteText(ctx, updatePointer)
		}
	}
}

func (a *App) RouteText(ctx context.Context, update *tgbotapi.Update) {
	_, _, state, ok := handler.UserSessionsMap.Get(update.Message.Chat.ID)
	if !ok {
		logger.Error("Err no value in map usersessions", "app-RouteText()", nil)
		a.Handler.UnknownCommand(ctx, update.Message.Chat.ID, "🔴 No value in map (user sessions), chat_id = "+strconv.Itoa(int(update.Message.Chat.ID)))
	}
	switch state {
	case "create_tournament_msg":
		a.Handler.CreateTournament(ctx, update)
	case "create_matches_msg":
		a.Handler.CreateMatches(ctx, update)
	case "add_results_msg":
		a.Handler.AddMatchResults(ctx, update)
	}
}
