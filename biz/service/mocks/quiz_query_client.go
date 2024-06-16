package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockQuizQueryClient struct {
	mock.Mock
}

func (_m *MockQuizQueryClient) GetParticipantsUserIDs(ctx context.Context, quizId string) ([]string, string, error) {
	ret := _m.Called(ctx, quizId)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(quizId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(quizId)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(quizId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}



