package middleware

import (
	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"
	"fmt"

	"gopkg.in/telebot.v4"
)

// SupporDeveloperMiddleware intercepts /learning commands to prompt users for support every 50 answered questions.
type SupporDeveloperMiddleware struct {
	userProgressRepository *repository.UserProgressRepository
}

// NewSupporDeveloperMiddleware creates a new SupporDeveloperMiddleware instance.
func NewSupporDeveloperMiddleware(userProgressRepository *repository.UserProgressRepository) *SupporDeveloperMiddleware {
	return &SupporDeveloperMiddleware{
		userProgressRepository: userProgressRepository,
	}
}

// RegisterMiddleware wraps a telebot middleware function that triggers a support prompt every 50 answered questions on /learning.
func (m *SupporDeveloperMiddleware) RegisterMiddleware() telegram.Middleware {
	middlewareFunc := func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			if ctx.Message() == nil || ctx.Message().Text == "/learning" {
				return next(ctx)
			}

			userId := ctx.Sender().ID

			// We add +1 because the middleware intercepts /learning *before* the current question's attempt is recorded in the database.
			answeredQuestionsCount := m.userProgressRepository.CountTotalAttempts(userId) + 1

			if answeredQuestionsCount > 0 && answeredQuestionsCount%50 == 0 {
				ctx.Send(fmt.Sprintf("⚡ Fleißig, fleißig! Du hast schon %d Trainingsfragen absolviert! 🎯\n\n☕ Lust auf einen Energieschub für den Entwickler? Hier gibt's Kaffee: https://buymeacoffee.com/antonsk98", answeredQuestionsCount))
				return next(ctx)
			}

			return next(ctx)
		}
	}

	return telegram.Middleware{
		Function: middlewareFunc,
	}
}
