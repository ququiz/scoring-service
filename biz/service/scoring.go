package service

import (
	"context"
	"ququiz/lintang/scoring-service/biz/dal/rabbitmq"
	"ququiz/lintang/scoring-service/biz/domain"

	"go.uber.org/zap"
)

type QuizQueryClient interface {
	GetParticipantsUserIDs(ctx context.Context, quizID string) ([]string, string, error)
}

type QuizMongoRepo interface {
	UpdateParticipantScore(ctx context.Context, quizID, participantUserID string, finalScore uint64) error
}

type AuthGRPClient interface {
	GetUsersByIds(ctx context.Context, userIDs []string) ([]domain.User, error)
}

type ScoringProducer interface {
	SendQuizRecap(ctx context.Context, quizRecapMessage domain.QuizRecapMessage) error
}

type QuizQueryProducer interface {
	SendDeleteCache(ctx context.Context, deleteCacheMessage domain.DeleteCacheMessage) error
}

type ScoringService struct {
	leaderboardRedis rabbitmq.LeaderboardRedis
	quizQueryClient  QuizQueryClient
	qMogoRepo        QuizMongoRepo
	authClient       AuthGRPClient
	sProducer        ScoringProducer
	qProducer        QuizQueryProducer
}

func NewScoringService(l rabbitmq.LeaderboardRedis, q QuizQueryClient, qr QuizMongoRepo,
	auth AuthGRPClient, sP ScoringProducer, qProducer QuizQueryProducer) *ScoringService {
	return &ScoringService{l, q, qr, auth, sP, qProducer}
}

func (s *ScoringService) GetLeaderboard(ctx context.Context, quizID string) ([]domain.LeaderBoard, error) {
	l, err := s.leaderboardRedis.GetTopLeaderBoard(ctx, quizID)
	if err != nil {
		zap.L().Error("s.leaderboardRedis.GetTopLeaderBoard (GetLeaderboard) (ScoringService)", zap.Error(err))
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

type ParticipantScore struct {
	Score    uint64
	Username string
	Rank     uint64
	Email    string
}

func (s *ScoringService) RecapQuiz(ctx context.Context, quizID string) error {
	l, err := s.leaderboardRedis.GetTopLeaderBoard(ctx, quizID)
	if err != nil {
		zap.L().Error("s.leaderboardRedis.GetTopLeaderBoard", zap.Error(err))
		return err
	}
	participantsIDs, quizName, err := s.quizQueryClient.GetParticipantsUserIDs(ctx, quizID)
	if err != nil {
		zap.L().Debug("participants from quiz query service: ", zap.Strings("participaantsID", participantsIDs))
		zap.L().Error(" s.quizQueryClient.GetParticipantsUserIDs (RecapQuiz) (ScoringService)", zap.Error(err))
		return err
	}

	var lMap map[string]uint64 = make(map[string]uint64)
	for i := 0; i < len(l); i++ {
		lMap[l[i].Username] = l[i].Score
	}

	participantsDetails, err := s.authClient.GetUsersByIds(ctx, participantsIDs)
	if err != nil {
		zap.L().Error("s.authClient.GetUsersByIds (RecapQuiz) (ScoringService) ", zap.Error(err))
		return err
	}

	var participantsScoreMap map[string]ParticipantScore = make(map[string]ParticipantScore)
	for i := 0; i < len(l); i++ {

		participant := participantsDetails[i]
		zap.L().Debug("participant from grpc: ", zap.String("userEmail", participant.Email))

		participantsScoreMap[participant.ID] = ParticipantScore{
			Score:    lMap[participant.Username],
			Username: participant.Username,
			Email:    participant.Email,
			Rank:     uint64(i + 1),
		}
	}

	for i := 0; i < len(participantsIDs); i++ {
		err := s.qMogoRepo.UpdateParticipantScore(ctx, quizID, participantsIDs[i], participantsScoreMap[participantsIDs[i]].Score)
		if err != nil {
			zap.L().Error(" s.qMogoRepo.UpdateParticipantScore (RecapQuiz) (ScoringService)", zap.Error(err))
			return err
		}
	}

	quizRecapMsg := domain.QuizRecapMessage{
		UserEmails: []string{},
		Leaderboard: domain.LeaderboardQuizRecap{
			QuizID:   quizID,
			QuizName: quizName,
		},
	}

	leaderboards := []domain.UserRanks{}

	for _, participant := range participantsScoreMap {
		zap.L().Debug("user email: ", zap.String("userEmail", participant.Email))

		leaderboards = append(leaderboards, domain.UserRanks{
			Email:    participant.Email,
			Rank:     participant.Rank,
			Score:    participant.Score,
			Username: participant.Username,
		})
		quizRecapMsg.UserEmails = append(quizRecapMsg.UserEmails, participant.Email)

	}
	quizRecapMsg.Leaderboard.Leaderboards = leaderboards

	err = s.sProducer.SendQuizRecap(ctx, quizRecapMsg)
	if err != nil {
		zap.L().Error("s.sProducer.SendQuizRecap (RecapQuiz) (ScoringService)", zap.Error(err))
		return err
	}

	err = s.qProducer.SendDeleteCache(ctx, domain.DeleteCacheMessage{QuizID: quizID})
	if err != nil {
		zap.L().Error("s.qProducer.SendDeleteCache (RecapQuiz) (ScoringService)", zap.Error(err))
		return err
	}
	return nil
}
