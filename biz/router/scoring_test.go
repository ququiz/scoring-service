package router_test

import (
	"encoding/json"
	"errors"
	"ququiz/lintang/scoring-service/biz/domain"
	"ququiz/lintang/scoring-service/biz/router"
	"ququiz/lintang/scoring-service/biz/router/mocks"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type getLeaderboardRes struct {
	Leaderboard []domain.LeaderBoard `json:"leaderboard"`
}

func TestGetLeaderboard(t *testing.T) {
	var mockLeaderboard []domain.LeaderBoard
	mockUseCase := new(mocks.MockScoringService)

	for i := 0; i < 10; i++ {
		var mockRank domain.LeaderBoard
		err := faker.FakeData(&mockRank)
		assert.NoError(t, err)
		mockLeaderboard = append(mockLeaderboard, mockRank)
	}

	var mockQuizID string = "666d8faaed25031b0d947430"

	handler := router.LeaderboardHandler{
		Svc: mockUseCase,
	}

	t.Run("success", func(t *testing.T) {
		h := server.Default()
		h.GET("/api/v1/scoring/:quizID/leaderboard", handler.GetLeaderboard)
		mockUseCase.On("GetLeaderboard", mock.Anything, mockQuizID).Return(mockLeaderboard, nil)

		w := ut.PerformRequest(h.Engine, "GET", "/api/v1/scoring/666d8faaed25031b0d947430/leaderboard", nil,
			ut.Header{"Connection", "close"})
		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode())

		mockLeaderboardExpected := getLeaderboardRes{mockLeaderboard}
		mockLeaderboardJSON, err := json.Marshal(mockLeaderboardExpected)

		assert.NoError(t, err)
		assert.Equal(t, string(mockLeaderboardJSON), string(resp.Body()))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Bad Request", func(t *testing.T) {
		h := server.Default()
		h.GET("/api/v1/scoring/:quizID/leaderboard", handler.GetLeaderboard)
		mockUseCase.On("GetLeaderboard", mock.Anything, mockQuizID).Return(mockLeaderboard, nil)

		w := ut.PerformRequest(h.Engine, "GET", "/api/v1/scoring/asdas/leaderboard", nil,
			ut.Header{"Connection", "close"})
		resp := w.Result()
		assert.Equal(t, 400, resp.StatusCode())
		mockUseCase.AssertExpectations(t)
	})

}

func TestGetLeaderboardFailed(t *testing.T) {
	var mockLeaderboard []domain.LeaderBoard
	mockUseCase := new(mocks.MockScoringService)

	for i := 0; i < 10; i++ {
		var mockRank domain.LeaderBoard
		err := faker.FakeData(&mockRank)
		assert.NoError(t, err)
		mockLeaderboard = append(mockLeaderboard, mockRank)
	}

	var mockQuizID string = "666d8faaed25031b0d947430"

	handler := router.LeaderboardHandler{
		Svc: mockUseCase,
	}

	t.Run("Internal Server Error", func(t *testing.T) {
		h := server.Default()
		h.GET("/api/v1/scoring/:quizID/leaderboard", handler.GetLeaderboard)
		mockUseCase.On("GetLeaderboard", mock.Anything, mockQuizID).Return([]domain.LeaderBoard{}, domain.WrapErrorf(errors.New(""), domain.ErrInternalServerError, domain.MessageInternalServerError))

		w := ut.PerformRequest(h.Engine, "GET", "/api/v1/scoring/666d8faaed25031b0d947430/leaderboard", nil,
			ut.Header{"Connection", "close"})

		resp := w.Result()

		assert.Equal(t, 500, resp.StatusCode())
		mockUseCase.AssertExpectations(t)
	})
}

func TestRecapQuiz(t *testing.T) {
	mockUseCase := new(mocks.MockScoringService)
	var mockQuizID string = "666d8faaed25031b0d947430"
	mockUseCase.On("RecapQuiz", mock.Anything, mockQuizID).Return(nil)
	handler := router.LeaderboardHandler{
		Svc: mockUseCase,
	}

	t.Run("success", func(t *testing.T) {
		h := server.Default()
		h.POST("/scoring-internal/recap/:quizID", handler.RecapQuiz)
		w := ut.PerformRequest(h.Engine, "POST", "/scoring-internal/recap/666d8faaed25031b0d947430", nil,
			ut.Header{"Connection", "close"})
		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode())
		mockUseCase.AssertExpectations(t)
	})
	t.Run("Bad Request", func(t *testing.T) {
		h := server.Default()
		h.POST("/scoring-internal/recap/:quizID", handler.RecapQuiz)
		w := ut.PerformRequest(h.Engine, "POST", "/scoring-internal/recap/asd", nil,
			ut.Header{"Connection", "close"})
		resp := w.Result()
		assert.Equal(t, 400, resp.StatusCode())
		mockUseCase.AssertExpectations(t)
	})
}

func TestRecapQuizInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.MockScoringService)
	var mockQuizID string = "666d8faaed25031b0d947430"
	mockUseCase.On("RecapQuiz", mock.Anything, mockQuizID).Return(domain.WrapErrorf(errors.New(""), domain.ErrInternalServerError, domain.MessageInternalServerError))
	handler := router.LeaderboardHandler{
		Svc: mockUseCase,
	}

	t.Run("internal server error when quiz recap", func(t *testing.T) {
		h := server.Default()
		h.POST("/scoring-internal/recap/:quizID", handler.RecapQuiz)
		w := ut.PerformRequest(h.Engine, "POST", "/scoring-internal/recap/666d8faaed25031b0d947430", nil,
			ut.Header{"Connection", "close"})
		resp := w.Result()
		assert.Equal(t, 500, resp.StatusCode())
		mockUseCase.AssertExpectations(t)
	})
}


