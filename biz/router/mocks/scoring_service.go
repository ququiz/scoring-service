package mocks

import (
	"context"
	"ququiz/lintang/scoring-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockScoringService struct {
	mock.Mock
}

func (_m *MockScoringService) GetLeaderboard(ctx context.Context, quizID string) ([]domain.LeaderBoard, error) {
	ret := _m.Called(ctx, quizID)

	var r0 []domain.LeaderBoard
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.LeaderBoard); ok {
		r0 = rf(ctx, quizID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.LeaderBoard)
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

func (_m *MockScoringService) RecapQuiz(ctx context.Context, quizID string) error {
	ret := _m.Called(ctx, quizID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, quizID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
