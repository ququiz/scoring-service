package mocks

import (
	"context"
	"ququiz/lintang/scoring-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockScoringProducer struct {
	mock.Mock
}

func (_m *MockScoringProducer) SendQuizRecap(ctx context.Context, quizRecapMessage domain.QuizRecapMessage) error {
	ret := _m.Called(ctx, quizRecapMessage)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.QuizRecapMessage) error); ok {
		r0 = rf(ctx, quizRecapMessage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
