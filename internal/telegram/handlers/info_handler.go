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
				"🪪 Willkommen im Trainingslager, %s! 🎯\n\n"+
					"Bereit, den deutschen Einbürgerungstest im Sturm zu erobern? 🇩🇪 Dieser Bot ist dein persönlicher Coach mit allen offiziellen BAMF-Fragen.\n\n"+
					"⚡ Tippe auf /learning, um deine erste Mission zu starten, Punkte zu sammeln und dem heiligen Pass näher zu kommen!\n\n"+
					"--- \n\n"+
					"☕ Dir gefällt das Abenteuer? Unterstütze die Entwicklung und spendiere dem Entwickler einen Kaffee: https://buymeacoffee.com/antonsk98",
				name,
			)

			return c.Send(message)
		},
	}
}
