package service_test

import (
	"context"
	"errors"
	"ququiz/lintang/scoring-service/biz/domain"
	"ququiz/lintang/scoring-service/biz/service"
	"ququiz/lintang/scoring-service/biz/service/mocks"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLeaderboard(t *testing.T) {
	mockLeaderBoard := new(mocks.MockLeaderboardRedis)
	mockQuizQueryClient := new(mocks.MockQuizQueryClient)
	mockAuthClient := new(mocks.MockAuthGRPCClient)
	mockQuizMongoRepo := new(mocks.MockQuizMongoRepo)
	mockScoringProducer := new(mocks.MockScoringProducer)
	mockQuizQueryProducer := new(mocks.MockQuizQueryProducer)

	var leaderboards []domain.RedisLeaderBoard
	for i := 0; i < 10; i++ {
		var leaderboard domain.RedisLeaderBoard
		err := faker.FakeData(&leaderboard)
		assert.NoError(t, err)
		leaderboards = append(leaderboards, leaderboard)
	}

	mockQuizID := "666d8faaed25031b0d947430"

	t.Run("success GetLeaderboard", func(t *testing.T) {
		mockLeaderBoard.On("GetTopLeaderBoard", mock.Anything, mockQuizID).Return(leaderboards, nil)
		service := service.NewScoringService(mockLeaderBoard, mockQuizQueryClient, mockQuizMongoRepo, mockAuthClient, mockScoringProducer, mockQuizQueryProducer)
		leaderboardsRes, err := service.GetLeaderboard(context.TODO(), mockQuizID)
		assert.NoError(t, err)
		assert.NotEmpty(t, leaderboardsRes)
		assert.Equal(t, len(leaderboards), len(leaderboardsRes))
		mockLeaderBoard.AssertExpectations(t)
	})

}

func TestGetLeaderBoardNotFound(t *testing.T) {
	mockLeaderBoard := new(mocks.MockLeaderboardRedis)
	mockQuizQueryClient := new(mocks.MockQuizQueryClient)
	mockAuthClient := new(mocks.MockAuthGRPCClient)
	mockQuizMongoRepo := new(mocks.MockQuizMongoRepo)
	mockScoringProducer := new(mocks.MockScoringProducer)
	mockQuizQueryProducer := new(mocks.MockQuizQueryProducer)

	var leaderboards []domain.LeaderBoard
	for i := 0; i < 10; i++ {
		var leaderboard domain.LeaderBoard
		err := faker.FakeData(&leaderboard)
		assert.NoError(t, err)
		leaderboards = append(leaderboards, leaderboard)
	}

	mockQuizID := "666d8faaed25031b0d947430"

	t.Run("not found GetLeaderboard", func(t *testing.T) {
		mockLeaderBoard.On("GetTopLeaderBoard", mock.Anything, mockQuizID).Return([]domain.RedisLeaderBoard{}, domain.WrapErrorf(errors.New(""), domain.ErrNotFound, "leaderboard not found"))
		service := service.NewScoringService(mockLeaderBoard, mockQuizQueryClient, mockQuizMongoRepo, mockAuthClient, mockScoringProducer, mockQuizQueryProducer)
		_, err := service.GetLeaderboard(context.TODO(), mockQuizID)
		assert.Error(t, err)
		mockLeaderBoard.AssertExpectations(t)
	})

}
