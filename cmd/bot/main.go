package main

import (
	"RushBananaBet/internal/app"
	"RushBananaBet/internal/handler"
	user "RushBananaBet/internal/model"
	"RushBananaBet/internal/repository"
	"RushBananaBet/internal/service"
	"RushBananaBet/pkg/db"
	"RushBananaBet/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"
)

var initData struct {
	Admins          []string `yaml:"admins"`
	BotToken        string   `yaml:"botToken"`
	LogLevel        int      `yaml:"log-level"`
	Production      bool     `yaml:"production"`
	User_db         string   `yaml:"user_db"`
	PasswordUser_db string   `yaml:"passwordUser_db"`
	Name_db         string   `yaml:"name_db"`
	Port_db         string   `yaml:"port_db"`
}

func init() {
	data, err := os.ReadFile("../../configs/config.yml")
	if err != nil {
		logger.Fatal("Cant read config", "main-init()", err)
	}

	err = yaml.Unmarshal(data, &initData)
	if err != nil {
		logger.Fatal("Cant unmarshal data from config", "main-init()", err)
	}

	// Set admins
	user.Admins = initData.Admins

	// Init logger
	logger.InitLogger(initData.LogLevel, initData.Production)
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	dbURL := "postgres://" +
		initData.User_db + ":" +
		initData.PasswordUser_db + "@localhost:" +
		initData.Port_db + "/" + initData.Name_db +
		"?sslmode=disable"

	db, err := db.NewPostgresPool(dbURL)
	if err != nil {
		logger.Fatal("Error create new pgxpool", "main-main()", err)
		return
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	bot := app.NewApp(initData.BotToken, *handler)
	bot.Start(stop)
}
