package domain

type LeaderBoard struct {
	Username string `json:"username"`
	Position uint64 `json:"position"`
	Score uint64 `json:"score"`
}


type RedisLeaderBoard struct {
	UserID string `json:"user_id"`
	Position uint64 `json:"position"`
	Score uint64 `json:"score"`
}
