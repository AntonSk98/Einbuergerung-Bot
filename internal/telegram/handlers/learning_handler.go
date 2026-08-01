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
				time.Sleep(2 * time.Second)
				_ = ctx.Delete()
			}()

			questionId, selectedOption, err := h.parseCallbackPayload(ctx)
			if err != nil {
				return ctx.Respond(&telebot.CallbackResponse{Text: err.Error()})
			}

			question, err := h.questionRepository.FindQuestionById(questionId)
			if err != nil {
				log.Printf("Failed to find question %d: %v", questionId, err)
				return ctx.Respond(&telebot.CallbackResponse{Text: "Frage konnte nicht gefunden werden."})
			}

			feedbackText, err := h.evaluateAnswer(userId, questionId, selectedOption, question.CorrectAnswer)
			if err != nil {
				log.Printf("Failed to update progress for user %d: %v", userId, err)
			}

			if err := ctx.Respond(&telebot.CallbackResponse{Text: feedbackText}); err != nil {
				log.Printf("Failed to respond to callback for user %d: %v", userId, err)
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
	return ctx.Send("Es ist ein Fehler aufgetreten. Bitte versuche es später noch einmal.")
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
func (h *LearningHandler) evaluateAnswer(userId int64, questionId int, selectedOption, correctAnswer string) (string, error) {
	if selectedOption == correctAnswer {
		err := h.userProgressRepository.HandleCorrectAnswer(userId, questionId)
		return "Richtig! :)", err
	}

	err := h.userProgressRepository.HandleWrongAnswer(userId, questionId)
	return "Falsch! :(", err
}
