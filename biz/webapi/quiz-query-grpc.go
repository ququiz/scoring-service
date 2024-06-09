package webapi

import (
	"context"
	"ququiz/lintang/scoring-service/config"
	"ququiz/lintang/scoring-service/kitex_gen/quiz-query-service/pb"
	"ququiz/lintang/scoring-service/kitex_gen/quiz-query-service/pb/quizqueryservice"
	"time"

	"github.com/cloudwego/kitex/client"
	"go.uber.org/zap"
)

type QuizQueryClient struct {
	service quizqueryservice.Client
}

func NewQuizQueryClient(cfg *config.Config) *QuizQueryClient {
	c, err := quizqueryservice.NewClient("quizQueryGRPCService", client.WithHostPorts(cfg.QueryQueryGRPC))
	if err != nil {
		zap.L().Fatal(" quizqueryservice.NewClient")
	}
	return &QuizQueryClient{c}
}

func (q *QuizQueryClient) GetParticipantsUserIDs(ctx context.Context, quizID string) ([]string, string, error) {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetQuizParticipantsReq{
		QuizId: quizID,
	}

	
	res, err := q.service.GetQuizParticipants(grpcCtx, req)
	if err != nil {
		zap.L().Error("q.service.GetQuizParticipants (GetParticipantsUserIDs) (QuizQueryClient)")
		return []string{}, "", err 
	}
	return res.UserIds, res.QuizName, nil 
}
