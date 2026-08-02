package main

import (
	"log"

	"einbuergerung-bot/internal/config"
	"einbuergerung-bot/internal/database"
	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"
	"einbuergerung-bot/internal/telegram/handlers"
	"einbuergerung-bot/internal/telegram/middleware"

	"github.com/jmoiron/sqlx"
)

func main() {

	config := loadConfig()
	databse := initDatabase(config)

	questionRepository := repository.NewQuestionRepository(databse)
	userRepository := repository.NewUserRepository(databse)
	userProgressRepository := repository.NewUserProgressRepository(databse)
	resetProgressRepository := repository.NewResetProgressRepository(databse)

	bot := newBot(config)

	infoHandler := handlers.NewInfoHandler()
	learningHandler := handlers.NewLearningHandler(userRepository, questionRepository, userProgressRepository)
	selectStateHandler := handlers.NewSelectFederalStateCallbackHandler(userRepository)
	progressHandler := handlers.NewProgressHandler(userProgressRepository)
	resetFederalStateHandler := handlers.NewResetFederalStateHandler(userRepository)
	resetProgressRepositoryHandler := handlers.NewResetProgressHandler(resetProgressRepository)

	authorizedMiddleware := middleware.NewAuthorizedUserMiddleware(config.AuthorizedUserIds)
	federalStateMiddleware := middleware.NewFederalStateSelectedMiddleware(userRepository)
	supportDeveloperMiddleware := middleware.NewSupporDeveloperMiddleware(userProgressRepository)

	middlewares := []telegram.Middleware{
		authorizedMiddleware.RegisterMiddleware(),
		federalStateMiddleware.RegisterMiddleware(),
		supportDeveloperMiddleware.RegisterMiddleware(),
	}

	handlersList := []telegram.Handler{
		infoHandler.RegisterCommand(),
		learningHandler.RegisterCommand(),
		learningHandler.RegisterCallback(),
		selectStateHandler.RegisterCallback(),
		progressHandler.RegisterCommand(),
		resetFederalStateHandler.RegisterCommand(),
		resetProgressRepositoryHandler.RegisterCommand(),
		resetProgressRepositoryHandler.RegisterCallbackStepOne(),
		resetProgressRepositoryHandler.RegisterCallbackStepTwo(),
	}

	bot.RegisterMiddleware(middlewares)
	bot.RegisterHandlers(handlersList)

	startBot(*bot)
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return cfg
}

func initDatabase(cfg *config.Config) *sqlx.DB {
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return db
}

func newBot(cfg *config.Config) *telegram.Bot {
	bot, err := telegram.NewBot(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	return bot
}

func startBot(bot telegram.Bot) {
	err := bot.Start()
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
