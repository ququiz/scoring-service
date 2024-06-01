package service

import (
	"context"
	"ququiz/lintang/scoring-service/biz/dal/domain"
	"ququiz/lintang/scoring-service/biz/dal/rabbitmq"

	"go.uber.org/zap"
)

type ScoringService struct {
	leaderboardRedis rabbitmq.LeaderboardRedis
}

func NewScoringService(l rabbitmq.LeaderboardRedis) *ScoringService {
	return &ScoringService{l}
}

func (s *ScoringService) GetLeaderboard(ctx context.Context, quizID string) ([]domain.LeaderBoard, error) {
	l, err := s.leaderboardRedis.GetTopLeaderBoard(ctx, quizID)
	if err != nil {
		zap.L().Error("s.leaderboardRedis.GetTopLeaderBoard", zap.Error(err))
		return []domain.LeaderBoard{}, err
	}
	var res []domain.LeaderBoard
	for i := 0; i < len(l); i++ {
		res = append(res, domain.LeaderBoard{
			Username: l[i].Username,
			Position: l[i].Position,
			Score:    l[i].Score,
		})
	}
	return res, nil
}
