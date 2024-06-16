package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockQuizMongoRepo struct {
	mock.Mock
}

func (_m *MockQuizMongoRepo) UpdateParticipantScore(ctx context.Context, quizID, participantUserID string, finalScore uint64) error {
	ret := _m.Called(ctx, quizID, participantUserID, finalScore)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uint64) error); ok {
		r0 = rf(ctx, quizID, participantUserID, finalScore)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


