// Package telegram handles Telegram bot initialization, routing, middlewares, and message management.
package telegram

import (
	"einbuergerung-bot/internal/config"
	"log"
	"time"

	"gopkg.in/telebot.v4"
)

// Bot represents the Telegram bot.
type Bot struct {
	bot *telebot.Bot
}

// Middleware wraps a telebot middleware function.
type Middleware struct {
	Function telebot.MiddlewareFunc
}

// Handler maps a Telegram endpoint to a specific handler function.
type Handler struct {
	Endpoint any
	Function telebot.HandlerFunc
}

// NewBot initializes and returns a new Telegram bot instance with a long poller.
func NewBot(cfg *config.Config) (*Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	return &Bot{bot: bot}, nil
}

// RegisterMiddleware registers a slice of custom middlewares to the bot.
func (b *Bot) RegisterMiddleware(middlewares []Middleware) {
	for _, middleware := range middlewares {
		b.bot.Use(middleware.Function)
	}
}

// RegisterHandlers registers a slice of handlers to their corresponding endpoints.
func (b *Bot) RegisterHandlers(handlers []Handler) {
	for _, handler := range handlers {
		b.bot.Handle(handler.Endpoint, handler.Function)
	}
}

// Handle registers a single handler function to a specific endpoint.
func (b *Bot) Handle(endpoint interface{}, h telebot.HandlerFunc) {
	b.bot.Handle(endpoint, h)
}

// Delete removes a given message from the chat.
func (b *Bot) Delete(msg *telebot.Message) error {
	return b.bot.Delete(msg)
}

// DeleteMessageAfterDelay waits for the specified delay before deleting a message.
func (b *Bot) DeleteMessageAfterDelay(msg *telebot.Message, delay time.Duration) {
	time.Sleep(delay)
	_ = b.Delete(msg)
}

// Start launches the Telegram bot polling loop.
func (b *Bot) Start() error {
	log.Println("Bot is running...")
	b.bot.Start()
	return nil
}
