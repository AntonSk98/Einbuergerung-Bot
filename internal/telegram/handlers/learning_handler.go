package handlers

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"einbuergerung-bot/internal/models"
	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"

	"gopkg.in/telebot.v4"
)

// LearningHandler manages learning sessions, question delivery, and answer evaluation.
type LearningHandler struct {
	userRepository         *repository.UserRepository
	questionRepository     *repository.QuestionRepository
	userProgressRepository *repository.UserProgressRepository
}

// NewLearningHandler creates a new instance of LearningHandler with the required repositories.
func NewLearningHandler(
	userRepository *repository.UserRepository,
	questionRepository *repository.QuestionRepository,
	userProgressRepository *repository.UserProgressRepository,
) *LearningHandler {
	return &LearningHandler{
		userRepository:         userRepository,
		questionRepository:     questionRepository,
		userProgressRepository: userProgressRepository,
	}
}

// RegisterCommand registers the /learning command handler.
func (h *LearningHandler) RegisterCommand() telegram.Handler {
	return telegram.Handler{
		Endpoint: "/learning",
		Function: h.handleLearningCommand,
	}
}

// RegisterCallback registers the inline keyboard callback handler for processing answered questions.
func (h *LearningHandler) RegisterCallback() telegram.Handler {
	return telegram.Handler{
		Endpoint: &telebot.Btn{Unique: "on_question_answered"},
		Function: h.handleAnswerCallback,
	}
}

// handleLearningCommand processes the /learning command and triggers the first question dispatch.
func (h *LearningHandler) handleLearningCommand(ctx telebot.Context) error {
	return h.sendQuestion(ctx, ctx.Sender().ID)
}

// handleAnswerCallback handles the user's inline button selection, evaluates the answer, manages hints, and loads the next question.
func (h *LearningHandler) handleAnswerCallback(ctx telebot.Context) error {
	userId := ctx.Sender().ID

	// Schedule auto-deletion of the answered question block
	h.scheduleMessageDeletion(ctx)

	questionId, selectedOption, err := h.parseCallbackPayload(ctx)
	if err != nil {
		return h.handleError(ctx, userId, err)
	}

	question, err := h.questionRepository.FindQuestionById(questionId)
	if err != nil {
		return h.handleError(ctx, userId, err)
	}

	answeredCorrectly, feedbackText, err := h.evaluateAnswer(userId, selectedOption, question)
	if err != nil {
		return h.handleError(ctx, userId, err)
	}

	if !answeredCorrectly {
		h.sendCorrectAnswerHint(ctx, question)
	}

	if err := ctx.Respond(&telebot.CallbackResponse{Text: feedbackText}); err != nil {
		return h.handleError(ctx, userId, err)
	}

	return h.sendQuestion(ctx, userId)
}

// sendQuestion fetches the user's next relevant test question, formats it (with photo support if applicable), and sends it to chat.
func (h *LearningHandler) sendQuestion(ctx telebot.Context, userId int64) error {
	user, err := h.userRepository.FindUserById(userId)
	if err != nil {
		return h.handleError(ctx, userId, err)
	}

	question, err := h.questionRepository.SelectQuestion(int64(user.ID), user.SelectedFederalState)
	if err != nil {
		return h.handleError(ctx, userId, err)
	}

	markup := h.buildOptionMarkup(question)
	questionText := h.formatQuestionText(question)

	if len(question.ImageData) > 0 {
		photo := &telebot.Photo{
			File:    telebot.FromReader(bytes.NewReader(question.ImageData)),
			Caption: questionText,
		}
		return ctx.Send(photo, markup, telebot.ModeMarkdown)
	}

	return ctx.Send(questionText, markup, telebot.ModeMarkdown)
}

// formatQuestionText structures the question text and all option mapping items cleanly into a single message string.
func (h *LearningHandler) formatQuestionText(q *models.Question) string {
	return fmt.Sprintf(
		"🗂️ Prüfungsfrage des Amtes\n\n"+
			"_%s_\n\n"+
			"🇦 %s\n\n"+
			"🇧 %s\n\n"+
			"🇨 %s\n\n"+
			"🇩 %s",
		q.QuestionText,
		q.OptionA,
		q.OptionB,
		q.OptionC,
		q.OptionD,
	)
}

// buildOptionMarkup constructs an inline keyboard layout containing selection buttons for options A, B, C, and D.
func (h *LearningHandler) buildOptionMarkup(question *models.Question) *telebot.ReplyMarkup {
	markup := &telebot.ReplyMarkup{}

	payload := func(opt string) string {
		return fmt.Sprintf("%d:%s", question.ID, opt)
	}

	btnA := markup.Data("🇦", "on_question_answered", payload("option_a"))
	btnB := markup.Data("🇧", "on_question_answered", payload("option_b"))
	btnC := markup.Data("🇨", "on_question_answered", payload("option_c"))
	btnD := markup.Data("🇩", "on_question_answered", payload("option_d"))

	markup.Inline(
		markup.Row(btnA, btnB),
		markup.Row(btnC, btnD),
	)

	return markup
}

// evaluateAnswer checks if the selected choice matches the correct answer, updates user progress, and returns feedback text.
func (h *LearningHandler) evaluateAnswer(userId int64, selectedOption string, question *models.Question) (bool, string, error) {
	if selectedOption == question.CorrectAnswer {
		err := h.userProgressRepository.HandleCorrectAnswer(userId, question.ID)
		return true, "✅ Richtig! +1 XP 😇", err
	}

	err := h.userProgressRepository.HandleWrongAnswer(userId, question.ID)
	return false, "❌ Falsch! -1 XP 😈", err
}

// parseCallbackPayload splits and validates the incoming callback data string into a question ID and selected option key.
func (h *LearningHandler) parseCallbackPayload(ctx telebot.Context) (int, string, error) {
	parts := strings.Split(ctx.Data(), ":")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("ungültige Antwortdaten")
	}

	questionId, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("ungültige Fragen-ID")
	}

	return questionId, parts[1], nil
}

// sendCorrectAnswerHint sends a temporary hint message revealing the correct answer if the user answered incorrectly.
func (h *LearningHandler) sendCorrectAnswerHint(ctx telebot.Context, question *models.Question) {
	hintMsg, err := ctx.Bot().Send(ctx.Recipient(), fmt.Sprintf("❌ Falscher Stempel! Richtig wäre:\n\n%s", question.GetCorrectAnswerText()), telebot.ModeMarkdown)

	if err == nil && hintMsg != nil {
		go func() {
			time.Sleep(30 * time.Second)
			_ = ctx.Bot().Delete(hintMsg)
		}()
	}
}

// scheduleMessageDeletion automatically deletes a given message after a 30-second delay.
func (h *LearningHandler) scheduleMessageDeletion(ctx telebot.Context) {
	go func() {
		time.Sleep(30 * time.Second)
		_ = ctx.Delete()
	}()
}

// handleError logs operational exceptions and responds to the user with a standardized error message.
func (h *LearningHandler) handleError(ctx telebot.Context, userId int64, err error) error {
	log.Printf("Failed for user %d: %v", userId, err)
	return ctx.Send("🚨 Alarm in der Zentrale! Etwas ist schiefgelaufen. Bitte lade deine Energie auf und versuche es gleich noch einmal! 🔄")
}
