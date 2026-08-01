package models

import "errors"

// Question represents a question in the database.
type Question struct {
	ID            int    `db:"id" json:"id"`
	FederalState  string `db:"federal_state" json:"federal_state"`
	QuestionText  string `db:"question_text" json:"question_text"`
	ImageData     []byte `db:"image_data" json:"image_data"`
	OptionA       string `db:"option_a" json:"option_a"`
	OptionB       string `db:"option_b" json:"option_b"`
	OptionC       string `db:"option_c" json:"option_c"`
	OptionD       string `db:"option_d" json:"option_d"`
	CorrectAnswer string `db:"correct_answer" json:"correct_answer"`
}

// UserProgress represents user progress in the database.
type UserProgress struct {
	UserID                 int    `db:"user_id"`
	QuestionID             int    `db:"question_id"`
	CorrectAnswerWeight    int    `db:"correct_answer_weight"`
	LastAnsweredQuestionAt string `db:"last_answered_question_at"`
}

// User represents a user in the database.
type User struct {
	ID                   int    `db:"telegram_id"`
	SelectedFederalState string `db:"selected_federal_state"`
}

const (
	StateBadenWuerttemberg     = "baden-wuerttemberg"
	StateBavaria               = "bavaria"
	StateBerlin                = "berlin"
	StateBrandenburg           = "brandenburg"
	StateBremen                = "bremen"
	StateHamburg               = "hamburg"
	StateHesse                 = "hesse"
	StateMecklenburgVorpommern = "mecklenburg-vorpommern"
	StateLowerSaxony           = "lower-saxony"
	StateNorthRhineWestphalia  = "north-rhine-westphalia"
	StateRhinelandPalatinate   = "rhineland-palatinate"
	StateSaarland              = "saarland"
	StateSaxony                = "saxony"
	StateSaxonyAnhalt          = "saxony-anhalt"
	StateSchleswigHolstein     = "schleswig-holstein"
	StateThuringia             = "thuringia"
)

// ResolveFederalStateByCode validates and returns the federal state code if it matches a known constant.
func ResolveFederalStateByCode(stateCode string) (string, error) {
	switch stateCode {
	case StateBadenWuerttemberg,
		StateBavaria,
		StateBerlin,
		StateBrandenburg,
		StateBremen,
		StateHamburg,
		StateHesse,
		StateMecklenburgVorpommern,
		StateLowerSaxony,
		StateNorthRhineWestphalia,
		StateRhinelandPalatinate,
		StateSaarland,
		StateSaxony,
		StateSaxonyAnhalt,
		StateSchleswigHolstein,
		StateThuringia:
		return stateCode, nil
	default:
		return "", errors.New("invalid federal state code")
	}
}

// GetCorrectAnswerText retrieves the actual text of the correct answer based on the stored option key (option_a, option_b, etc.).
func (q *Question) GetCorrectAnswerText() string {
	switch q.CorrectAnswer {
	case "option_a":
		return q.OptionA
	case "option_b":
		return q.OptionB
	case "option_c":
		return q.OptionC
	case "option_d":
		return q.OptionD
	default:
		return ""
	}
}
