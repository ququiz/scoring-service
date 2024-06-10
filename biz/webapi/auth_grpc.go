package webapi

import (
	"context"
	"ququiz/lintang/scoring-service/biz/domain"
	"ququiz/lintang/scoring-service/pb"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthClient struct {
	service pb.UsersServiceClient
}

func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	svc := pb.NewUsersServiceClient(cc)
	return &AuthClient{service: svc}
}

func (a *AuthClient) GetUsersByIds(ctx context.Context, userIDs []string) ([]domain.User, error) {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	zap.L().Debug("userIDs: ", zap.Strings("userIDs", userIDs))

	req := &pb.GetUserRequestByIds{
		Ids: userIDs,
	}

	res, err := a.service.GetUserByIds(grpcCtx, req)
	if err != nil {
		zap.L().Error("m.service.GetUserByIds  (GetUsersByIds) (UserGRPClient)", zap.Error(err))
		return []domain.User{}, err
	}

	var usernames []domain.User
	for i := 0; i < len(res.Users); i++ {
		zap.L().Debug("user from grpc: ", zap.String("userEmail", res.Users[i].Email))
		usernames = append(usernames, domain.User{
			ID:       res.Users[i].Id,
			Username: res.Users[i].Username,
			Email:    res.Users[i].Email,
			Fullname: res.Users[i].Username,
		})
	}

	return usernames, nil
}
