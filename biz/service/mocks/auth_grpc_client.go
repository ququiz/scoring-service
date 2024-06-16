package mocks

import (
	"context"
	"ququiz/lintang/scoring-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockAuthGRPCClient struct {
	mock.Mock
}

func (_m *MockAuthGRPCClient) GetUsersByIds(ctx context.Context, userIDs []string) ([]domain.User, error) {
	ret := _m.Called(ctx, userIDs)

	var r0 []domain.User
	if rf, ok := ret.Get(0).(func(context.Context, []string) []domain.User); ok {
		r0 = rf(ctx, userIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, userIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

