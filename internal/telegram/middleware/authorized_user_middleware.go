package middleware

import (
	"einbuergerung-bot/internal/telegram"
	"slices"

	"gopkg.in/telebot.v4"
)

// AuthorizedUserMiddleware restricts bot access to a predefined whitelist of user identifiers.
type AuthorizedUserMiddleware struct {
	authorizedUserIdentifiers []int64
}

// NewAuthorizedUserMiddleware creates a new instance of AuthorizedUserMiddleware with the given whitelist.
func NewAuthorizedUserMiddleware(authorizedUserIdentifiers []int64) *AuthorizedUserMiddleware {
	return &AuthorizedUserMiddleware{
		authorizedUserIdentifiers: authorizedUserIdentifiers,
	}
}

// RegisterMiddleware wraps a telebot middleware function that validates whether the sender is authorized.
func (m *AuthorizedUserMiddleware) RegisterMiddleware() telegram.Middleware {
	middlewareFunc := func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			currentUserID := c.Sender().ID

			isAllowed := slices.Contains(m.authorizedUserIdentifiers, currentUserID)

			if !isAllowed {
				return c.Send("⛔ You are not authorized to use this bot.")
			}

			return next(c)
		}
	}

	return telegram.Middleware{
		Function: middlewareFunc,
	}
}
