package domain

type CorrectAnswer struct {
	Weight   uint64 `json:"weight"`
	UserID   string `json:"user_id"`
	Username string `json:"user_name"`
	QuizID   string `json:"quiz_id"`
}
