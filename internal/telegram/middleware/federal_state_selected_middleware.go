// Package middleware provides custom request interceptors and access control checks for Telegram bot commands.
package middleware

import (
	"einbuergerung-bot/internal/models"
	"einbuergerung-bot/internal/repository"
	"einbuergerung-bot/internal/telegram"

	"gopkg.in/telebot.v4"
)

// FederalStateSelectedMiddleware verifies whether the user has chosen a federal state before allowing learning commands.
type FederalStateSelectedMiddleware struct {
	userRepository *repository.UserRepository
}

// NewFederalStateSelectedMiddleware creates a new instance of FederalStateSelectedMiddleware without caching.
func NewFederalStateSelectedMiddleware(userRepository *repository.UserRepository) *FederalStateSelectedMiddleware {
	return &FederalStateSelectedMiddleware{
		userRepository: userRepository,
	}
}

// RegisterMiddleware wraps a telebot middleware function that intercepts /learning commands if no federal state is set.
func (m *FederalStateSelectedMiddleware) RegisterMiddleware() telegram.Middleware {
	middlewareFunc := func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			if ctx.Message() == nil || ctx.Message().Text != "/learning" {
				return next(ctx)
			}

			userId := ctx.Sender().ID

			if m.userRepository.FederalStateSelected(userId) {
				return next(ctx)
			}

			selectFederalStateReplyMarkup := m.initReplyMarkup()

			return ctx.Send("🗺️ Achtung, Rekrut! Einige Missionen sind bundeslandspezifisch. Wähle zuerst dein Bundesland aus, um fortzufahren:", selectFederalStateReplyMarkup)
		}
	}

	return telegram.Middleware{
		Function: middlewareFunc,
	}
}

// initReplyMarkup builds an inline keyboard containing all German federal states for selection.
func (m *FederalStateSelectedMiddleware) initReplyMarkup() *telebot.ReplyMarkup {
	selectFederalStateReplyMarkup := &telebot.ReplyMarkup{}

	optionStateBadenWuerttemberg := selectFederalStateReplyMarkup.Data("Baden-Württemberg", "select_state", models.StateBadenWuerttemberg)
	optionStateBavaria := selectFederalStateReplyMarkup.Data("Bayern", "select_state", models.StateBavaria)
	optionStateBerlin := selectFederalStateReplyMarkup.Data("Berlin", "select_state", models.StateBerlin)
	optionStateBrandenburg := selectFederalStateReplyMarkup.Data("Brandenburg", "select_state", models.StateBrandenburg)
	optionStateBremen := selectFederalStateReplyMarkup.Data("Bremen", "select_state", models.StateBremen)
	optionStateHamburg := selectFederalStateReplyMarkup.Data("Hamburg", "select_state", models.StateHamburg)
	optionStateHesse := selectFederalStateReplyMarkup.Data("Hessen", "select_state", models.StateHesse)
	optionStateMecklenburgVorpommern := selectFederalStateReplyMarkup.Data("Mecklenburg-Vorpommern", "select_state", models.StateMecklenburgVorpommern)
	optionStateLowerSaxony := selectFederalStateReplyMarkup.Data("Niedersachsen", "select_state", models.StateLowerSaxony)
	optionStateNorthRhineWestphalia := selectFederalStateReplyMarkup.Data("Nordrhein-Westfalen", "select_state", models.StateNorthRhineWestphalia)
	optionStateRhinelandPalatinate := selectFederalStateReplyMarkup.Data("Rheinland-Pfalz", "select_state", models.StateRhinelandPalatinate)
	optionStateSaarland := selectFederalStateReplyMarkup.Data("Saarland", "select_state", models.StateSaarland)
	optionStateSaxony := selectFederalStateReplyMarkup.Data("Sachsen", "select_state", models.StateSaxony)
	optionStateSaxonyAnhalt := selectFederalStateReplyMarkup.Data("Sachsen-Anhalt", "select_state", models.StateSaxonyAnhalt)
	optionStateSchleswigHolstein := selectFederalStateReplyMarkup.Data("Schleswig-Holstein", "select_state", models.StateSchleswigHolstein)
	optionStateThuringia := selectFederalStateReplyMarkup.Data("Thüringen", "select_state", models.StateThuringia)

	selectFederalStateReplyMarkup.Inline(
		selectFederalStateReplyMarkup.Row(optionStateBadenWuerttemberg, optionStateBavaria),
		selectFederalStateReplyMarkup.Row(optionStateBerlin, optionStateBrandenburg),
		selectFederalStateReplyMarkup.Row(optionStateBremen, optionStateHamburg),
		selectFederalStateReplyMarkup.Row(optionStateHesse, optionStateMecklenburgVorpommern),
		selectFederalStateReplyMarkup.Row(optionStateLowerSaxony, optionStateNorthRhineWestphalia),
		selectFederalStateReplyMarkup.Row(optionStateRhinelandPalatinate, optionStateSaarland),
		selectFederalStateReplyMarkup.Row(optionStateSaxony, optionStateSaxonyAnhalt),
		selectFederalStateReplyMarkup.Row(optionStateSchleswigHolstein, optionStateThuringia),
	)

	return selectFederalStateReplyMarkup
}
