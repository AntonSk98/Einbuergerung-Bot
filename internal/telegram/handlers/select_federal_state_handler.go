package handlers

import (
	"einbuergerung-bot/internal/models"
	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"
	"time"

	"gopkg.in/telebot.v4"
)

// SelectFederalStateCallbackHandler handles inline button callbacks for selecting a federal state.
type SelectFederalStateCallbackHandler struct {
	userRepository *repository.UserRepository
}

// NewSelectFederalStateCallbackHandler creates a new instance of SelectFederalStateCallbackHandler.
func NewSelectFederalStateCallbackHandler(
	userRepository *repository.UserRepository,
) *SelectFederalStateCallbackHandler {
	return &SelectFederalStateCallbackHandler{
		userRepository: userRepository,
	}
}

// RegisterCallback registers the inline keyboard callback handler for processing federal state choices.
func (h *SelectFederalStateCallbackHandler) RegisterCallback() telegram.Handler {
	return telegram.Handler{
		Endpoint: &telebot.Btn{Unique: "select_state"},
		Function: func(c telebot.Context) error {
			payload := c.Callback().Data

			go func() {
				time.Sleep(1 * time.Second)
				_ = c.Delete()
			}()

			federalState, err := models.ResolveFederalStateByCode(payload)
			if err != nil {
				return c.Respond(&telebot.CallbackResponse{
					Text: "Ungültiges Bundesland. Bitte versuche es erneut.",
				})
			}

			userId := c.Sender().ID

			if err := h.userRepository.PersistFederalState(userId, federalState); err != nil {
				return c.Respond(&telebot.CallbackResponse{
					Text: "Fehler beim Speichern. Bitte versuche es erneut.",
				})
			}

			return c.Send("Danke! Deine Auswahl wurde gespeichert. Du kannst dich jetzt auf deinen Einbürgerungstest vorbereiten! Starte mit /learning!")
		},
	}
}
