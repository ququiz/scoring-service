package domain

type QuizRecapMessage struct {
	UserEmails []string `json:"user_emails"`
	Leaderboard LeaderboardQuizRecap `json:"leaderboard"`
}

type LeaderboardQuizRecap struct {
	QuizID string `json:"quiz_id"`
	QuizName string `json:"quiz_name"`
	Leaderboards []UserRanks `json:"leaderboards"`
}	

type UserRanks struct {
	Email string `json:"email"`
	Rank uint64 `json:"rank"`
	Score uint64 `json:"score"`
	Username string `json:"username"`
}


type DeleteCacheMessage struct {
	QuizID string `json:"quiz_id"`
}

