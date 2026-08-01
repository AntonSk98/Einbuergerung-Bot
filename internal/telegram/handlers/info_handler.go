package handlers

import (
	"einbuergerung-bot/internal/telegram"
	"fmt"

	"gopkg.in/telebot.v4"
)

// InfoHandler handles general informational and welcome commands.
type InfoHandler struct{}

// NewInfoHandler creates a new instance of InfoHandler.
func NewInfoHandler() *InfoHandler {
	return &InfoHandler{}
}

// RegisterCommand registers the /start command handler, greeting the user and introducing the bot.
func (h *InfoHandler) RegisterCommand() telegram.Handler {
	return telegram.Handler{
		Endpoint: "/start",
		Function: func(c telebot.Context) error {
			user := c.Sender()

			name := user.FirstName
			if name == "" {
				name = user.Username
			}

			message := fmt.Sprintf(
				"👋 Willkommen, %s!\n\n"+
					"Dieser Bot hilft dir dabei, den deutschen Einbürgerungstest erfolgreich zu bestehen. "+
					"Er enthält alle offiziellen Fragen des BAMF.\n\n"+
					"Verwende den Befehl /learning, um mit dem Lernen zu beginnen!",
				name,
			)

			return c.Send(message)
		},
	}
}
