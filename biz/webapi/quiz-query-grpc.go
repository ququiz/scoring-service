package webapi

import (
	"context"
	"fmt"
	"ququiz/lintang/scoring-service/biz/domain"
	"ququiz/lintang/scoring-service/pb"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// type QuizQueryClient struct {
// 	service quizqueryservice.Client
// }

// func NewQuizQueryClient(cfg *config.Config) *QuizQueryClient {
// 	c, err := quizqueryservice.NewClient("quizQueryGRPCService", client.WithHostPorts(cfg.QueryQueryGRPC))
// 	if err != nil {
// 		zap.L().Fatal(" quizqueryservice.NewClient", zap.Error(err))
// 	}

// 	return &QuizQueryClient{c}
// }

// func (q *QuizQueryClient) GetParticipantsUserIDs(ctx context.Context, quizID string) ([]string, string, error) {
// 	grpcCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	req := &pb.GetQuizParticipantsReq{
// 		QuizId: quizID,
// 	}

// 	res, err := q.service.GetQuizParticipants(grpcCtx, req)
// 	if err != nil {
// 		zap.L().Error("q.service.GetQuizParticipants (GetParticipantsUserIDs) (QuizQueryClient)", zap.Error(err))
// 		return []string{}, "", err
// 	}
// 	return res.UserIds, res.QuizName, nil
// }

type QuizQueryClient struct {
	service pb.QuizQueryServiceClient
}

func NewQuizQueryClient(cc *grpc.ClientConn) *QuizQueryClient {
	svc := pb.NewQuizQueryServiceClient(cc)

	return &QuizQueryClient{service: svc}
}

func (q *QuizQueryClient) GetParticipantsUserIDs(ctx context.Context, quizId string) ([]string, string, error) {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 4 * time.Second)
	defer cancel()

	zap.L().Debug(fmt.Sprintf(`quizID: %s`, quizId))
	req := &pb.GetQuizParticipantsReq{
		QuizId: quizId,
	}
	res, err := q.service.GetQuizParticipants(grpcCtx, req)
	if err != nil {
		zap.L().Error("q.service.GetQuizParticipants (GetParticipantsUserIDs) (QuizQueryClient)", zap.Error(err))
		return []string{}, "",domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	return res.UserIds, res.QuizName, nil
}
