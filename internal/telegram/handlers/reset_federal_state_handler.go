// Package handlers manages Telegram bot command and event handlers.
package handlers

import (
	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"

	"gopkg.in/telebot.v4"
)

// ResetFederalStateHandler handles the command to reset a user's chosen federal state.
type ResetFederalStateHandler struct {
	userRepository *repository.UserRepository
}

// NewResetFederalStateHandler creates a new instance of ResetFederalStateHandler.
func NewResetFederalStateHandler(userRepository *repository.UserRepository) *ResetFederalStateHandler {
	return &ResetFederalStateHandler{
		userRepository: userRepository,
	}
}

// RegisterCommand registers the /reset_federal_state endpoint with the Telegram bot.
func (h *ResetFederalStateHandler) RegisterCommand() telegram.Handler {
	return telegram.Handler{
		Endpoint: "/reset_federal_state",
		Function: func(c telebot.Context) error {
			userId := c.Sender().ID

			err := h.userRepository.ResetFederalState(userId)
			if err != nil {
				return c.Send("🚨 Oje! Der Amtsschimmel hat gehustet - da ist wohl digitaler Sand im Getriebe gelandet. Probier's noch mal!")
			}

			return c.Send("🧹 Tabula rasa! Dein Bundesland wurde erfolgreich aus den Akten getilgt. Keine Sorge, deine hart erkämpften XP bleiben sicher! Beim nächsten /learning wirst du einfach ein neues Bundesland wählen.")
		},
	}
}
