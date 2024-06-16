package mocks

import (
	"context"
	"ququiz/lintang/scoring-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockLeaderboardRedis struct {
	mock.Mock
}

func (_m *MockLeaderboardRedis) CalculateUserScore(ctx context.Context, weight uint64,
	userID string, quizID string, userName string) error {

	ret := _m.Called(ctx, weight, userID, quizID, userName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, string, string, string) error); ok {
		r0 = rf(ctx, weight, userID, quizID, userName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *MockLeaderboardRedis) GetTopLeaderBoard(ctx context.Context, quizID string) ([]domain.RedisLeaderBoard, error) {
	ret := _m.Called(ctx, quizID)

	var r0 []domain.RedisLeaderBoard
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.RedisLeaderBoard); ok {
		r0 = rf(ctx, quizID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.RedisLeaderBoard)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, quizID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

