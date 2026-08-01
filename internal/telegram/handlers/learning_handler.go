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

// RegisterCommand registers the /start-learning or /learning command handler.
func (h *LearningHandler) RegisterCommand() telegram.Handler {
	return telegram.Handler{
		Endpoint: "/learning",
		Function: func(ctx telebot.Context) error {
			return h.sendQuestion(ctx, ctx.Sender().ID)
		},
	}
}

// RegisterCallback registers the inline keyboard callback handler for processing answered questions.
func (h *LearningHandler) RegisterCallback() telegram.Handler {
	return telegram.Handler{
		Endpoint: &telebot.Btn{Unique: "on_question_answered"},
		Function: func(ctx telebot.Context) error {
			userId := ctx.Sender().ID

			go func() {
				time.Sleep(30 * time.Second)
				_ = ctx.Delete()
			}()

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
				// Show the correct answer with a game-like auto-deleting message
				msg, _ := ctx.Bot().Send(ctx.Recipient(), fmt.Sprintf("💡 Richtige Antwort wäre gewesen:\n\n%s", question.GetCorrectAnswerText()))
				if msg != nil {
					go func() {
						time.Sleep(30 * time.Second)
						_ = ctx.Bot().Delete(msg)
					}()
				}
			}
			if err := ctx.Respond(&telebot.CallbackResponse{Text: feedbackText}); err != nil {
				return h.handleError(ctx, userId, err)
			}

			return h.sendQuestion(ctx, userId)
		},
	}
}

// sendQuestion retrieves the next appropriate question for the user and sends it as text or photo with an inline keyboard.
func (h *LearningHandler) sendQuestion(ctx telebot.Context, userId int64) error {
	user, err := h.userRepository.FindUserById(userId)
	if err != nil {
		return h.handleError(ctx, userId, err)
	}

	randomQuestion, err := h.questionRepository.SelectQuestion(int64(user.ID), user.SelectedFederalState)
	if err != nil {
		return h.handleError(ctx, userId, err)
	}

	markup := h.initReplyMarkup(randomQuestion)

	if len(randomQuestion.ImageData) > 0 {
		photo := &telebot.Photo{
			File:    telebot.FromReader(bytes.NewReader(randomQuestion.ImageData)),
			Caption: randomQuestion.QuestionText,
		}
		return ctx.Send(photo, markup)
	}

	return ctx.Send(randomQuestion.QuestionText, markup)
}

// handleError logs operational errors and sends a generic error message to the user.
func (h *LearningHandler) handleError(ctx telebot.Context, userId int64, err error) error {
	log.Printf("Failed for user %d: %v", userId, err)
	return ctx.Send("🚨 Alarm in der Zentrale! Etwas ist schiefgelaufen. Bitte lade deine Energie auf und versuche es gleich noch einmal! 🔄")
}

// initReplyMarkup builds the inline keyboard markup containing the options for a question.
func (h *LearningHandler) initReplyMarkup(question *models.Question) *telebot.ReplyMarkup {
	questionReplyMarkup := &telebot.ReplyMarkup{}

	initCallbackPayload := func(option string) string {
		return fmt.Sprintf("%d:%s", question.ID, option)
	}

	optionA := questionReplyMarkup.Data(question.OptionA, "on_question_answered", initCallbackPayload("option_a"))
	optionB := questionReplyMarkup.Data(question.OptionB, "on_question_answered", initCallbackPayload("option_b"))
	optionC := questionReplyMarkup.Data(question.OptionC, "on_question_answered", initCallbackPayload("option_c"))
	optionD := questionReplyMarkup.Data(question.OptionD, "on_question_answered", initCallbackPayload("option_d"))

	questionReplyMarkup.Inline(
		questionReplyMarkup.Row(optionA),
		questionReplyMarkup.Row(optionB),
		questionReplyMarkup.Row(optionC),
		questionReplyMarkup.Row(optionD),
	)

	return questionReplyMarkup
}

// parseCallbackPayload extracts and validates the question ID and selected option from the callback data.
func (h *LearningHandler) parseCallbackPayload(ctx telebot.Context) (int, string, error) {
	payload := ctx.Data()
	parts := strings.Split(payload, ":")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("ungültige Antwortdaten")
	}

	questionId, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("ungültige Fragen-ID")
	}

	return questionId, parts[1], nil
}

// evaluateAnswer checks if the selected option matches the correct answer and updates user progress accordingly.
func (h *LearningHandler) evaluateAnswer(userId int64, selectedOption string, question *models.Question) (bool, string, error) {
	if selectedOption == question.CorrectAnswer {
		err := h.userProgressRepository.HandleCorrectAnswer(userId, question.ID)
		return true, "✅ Richtig! +1 XP 😇", err
	}

	err := h.userProgressRepository.HandleWrongAnswer(userId, question.ID)
	return false, "❌ Falsch! -1 XP 😈", err
}
