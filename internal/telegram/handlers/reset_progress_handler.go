package handlers

import (
	"strings"

	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"

	"gopkg.in/telebot.v4"
)

// ResetProgressHandler manages the multi-step confirmation flow to wipe user progress and start a new game.
type ResetProgressHandler struct {
	repository *repository.ResetProgressRepository
}

// NewResetProgressHandler creates a new instance of ResetProgressHandler.
func NewResetProgressHandler(repo *repository.ResetProgressRepository) *ResetProgressHandler {
	return &ResetProgressHandler{
		repository: repo,
	}
}

// RegisterCommand registers the /new-game command endpoint.
func (h *ResetProgressHandler) RegisterCommand() telegram.Handler {
	return telegram.Handler{
		Endpoint: "/new-game",
		Function: h.handleNewGameCommand,
	}
}

// RegisterCallbackStepOne registers the first inline keyboard callback for resetting user progress.
func (h *ResetProgressHandler) RegisterCallbackStepOne() telegram.Handler {
	return telegram.Handler{
		Endpoint: &telebot.Btn{Unique: "reset_step_1"},
		Function: h.handleStepOneCallback,
	}
}

// RegisterCallbackStepTwo registers the second inline keyboard callback for the final confirmation and execution of the reset.
func (h *ResetProgressHandler) RegisterCallbackStepTwo() telegram.Handler {
	return telegram.Handler{
		Endpoint: &telebot.Btn{Unique: "reset_step_2"},
		Function: h.handleStepTwoCallback,
	}
}

// handleNewGameCommand prompts the user with an initial confirmation warning before wiping progress.
func (h *ResetProgressHandler) handleNewGameCommand(c telebot.Context) error {
	markup := &telebot.ReplyMarkup{}

	btnYes := markup.Data("Ja, alles löschen 💀", "reset_step_1", "confirm_yes")
	btnNo := markup.Data("Nein, behalten! 😅", "reset_step_1", "confirm_no")

	markup.Inline(
		markup.Row(btnYes),
		markup.Row(btnNo),
	)

	message := "⚠️ Wollen wir wirklich von vorne beginnen?\n\nDas bedeutet: Deine hart erkämpften XP verpuffen im Nirgendwo. Bist du sicher?"
	return c.Send(message, markup, telebot.ModeMarkdown)
}

// handleStepOneCallback processes the first confirmation step, handling cancellation or advancing to the final warning.
func (h *ResetProgressHandler) handleStepOneCallback(c telebot.Context) error {
	if h.isActionDeclined(c.Data()) {
		_ = c.Respond(&telebot.CallbackResponse{Text: "Puh! Knapp an der Katastrophe vorbei! 😅"})
		return c.Edit("Puh! 💨 Deine Akten sind sicher im Amt geblieben. Weiter geht's mit dem Grind! 🎮🔥")
	}

	markup := h.buildFinalWarningMarkup()
	_ = c.Respond()

	warningMessage := "🚨 LETZTE WARNUNG! 🚨\n\nDer Amtsschimmel grinst böse. Es gibt kein Zurück mehr. Alle XP werden restlos geschreddert!\n\nWillst du das wirklich?"
	return c.Edit(warningMessage, markup, telebot.ModeMarkdown)
}

// handleStepTwoCallback executes the final reset of user progress or aborts based on the user's final choice.
func (h *ResetProgressHandler) handleStepTwoCallback(c telebot.Context) error {
	if h.isActionDeclined(c.Data()) {
		_ = c.Respond(&telebot.CallbackResponse{Text: "Puh... Schweiß von der Stirn wisch! 😅"})
		return c.Edit("Puh, wir sind gerade noch so dem bürokratischen Exil entkommen! Deine XP sind gerettet. 🛡️✨")
	}

	userId := c.Sender().ID
	if err := h.repository.ResetProgress(userId); err != nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: "Fehler beim Schreddern!"})
		return c.Edit("🚨 Upsi! Der Aktenschrank klemmt. Das Zurücksetzen ist fehlgeschlagen, du musst wohl weiter lernen!")
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: "Alles auf Null! 📉"})
	successMessage := "💥 Tabula Rasa!\n\nDein Profil wurde komplett vernichtet. Du bist wieder ganz unten auf dem harten Boden der Tatsachen aufgeschlagen. Nutze /learning, um deine neue Reise zu starten! 🎮🔄"
	return c.Edit(successMessage)
}

// isActionDeclined checks if the callback payload indicates that the user declined the reset action.
func (h *ResetProgressHandler) isActionDeclined(payloadData string) bool {
	payload := strings.Split(payloadData, ":")
	if len(payload) == 0 {
		return false
	}
	action := payload[len(payload)-1]
	return action == "confirm_no" || action == "final_no"
}

// buildFinalWarningMarkup creates the inline keyboard markup for the final confirmation step.
func (h *ResetProgressHandler) buildFinalWarningMarkup() *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}
	btnFinalYes := markup.Data("JA! ICH WILL LEIDEN! 🔥", "reset_step_2", "final_yes")
	btnFinalNo := markup.Data("STOPP! RETTE MICH! 🛑", "reset_step_2", "final_no")

	markup.Inline(
		markup.Row(btnFinalYes),
		markup.Row(btnFinalNo),
	)
	return markup
}
