package mocks

import (
	"context"
	"ququiz/lintang/scoring-service/biz/domain"

	"github.com/stretchr/testify/mock"
)

type MockQuizQueryProducer struct {
	mock.Mock
}

func (_M *MockQuizQueryProducer) SendDeleteCache(ctx context.Context, deleteCacheMessage domain.DeleteCacheMessage) error {
	ret := _M.Called(ctx, deleteCacheMessage)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.DeleteCacheMessage) error); ok {
		r0 = rf(ctx, deleteCacheMessage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

