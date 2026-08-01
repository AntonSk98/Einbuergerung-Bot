package handlers

import (
	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"
	"fmt"

	"gopkg.in/telebot.v4"
)

// ProgressHandler manages user progress tracking and profile statistics commands.
type ProgressHandler struct {
	userProgressRepository *repository.UserProgressRepository
}

// NewProgressHandler creates a new instance of ProgressHandler.
func NewProgressHandler(userProgressRepository *repository.UserProgressRepository) *ProgressHandler {
	return &ProgressHandler{
		userProgressRepository: userProgressRepository,
	}
}

// RegisterCommand registers the /progress command handler, displaying the user's net XP score and rank tier.
func (h *ProgressHandler) RegisterCommand() telegram.Handler {
	return telegram.Handler{
		Endpoint: "/progress",
		Function: func(c telebot.Context) error {
			userId := c.Sender().ID
			netScore := h.userProgressRepository.GetNetScore(userId)

			var title, vibe, roast string
			switch {
			case netScore < 0:
				title = "Vom Amt Verfluchter"
				vibe = "😈"
				roast = "Der Sachbearbeiter weint seinetwegen. Amtsschimmel-Alarm!"
			case netScore <= 99:
				title = "Asylsuchender im Amt"
				vibe = "👺"
				roast = "Das Faxgerät ist schneller als dein Wissensstand. Weiterlernen!"
			case netScore <= 249:
				title = "Stempel-Lehrling"
				vibe = "📄"
				roast = "Du hast immerhin schon mal eine Bescheinigung über die Bescheinigung gesehen."
			case netScore <= 399:
				title = "Wartezimmer-Profi"
				vibe = "⚖️"
				roast = "Du kennst die Öffnungszeiten vom Bürgeramt auswendig. Nicht schlecht!"
			case netScore <= 549:
				title = "Termin-König"
				vibe = "🛡️"
				roast = "Respekt! Du kriegst sogar innerhalb von zwei Wochen einen Bürgeramts-Termin."
			default:
				title = "Ehren-Deutscher mit Urkunde"
				vibe = "😇"
				roast = "Perfekt! Du bist bereit, eigenhändig Bratwurst zu wenden und Müll zu trennen."
			}

			message := fmt.Sprintf(
				"🪪 DEINE AKTE BEIM BÜRGERAMT 🪪\n\n"+
					"Aktueller Kontostand: `%d` XP\n"+
					"Dienstgrad: %s %s\n\n"+
					"💬 Gutachter-Kommentar:\n -> %s\n\n"+
					"🎯 Ziel für den heiligen Pass: `500+ XP`",
				netScore, title, vibe, roast,
			)

			return c.Send(message, telebot.ModeMarkdown)
		},
	}
}
