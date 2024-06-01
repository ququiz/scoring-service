package router

import (
	"context"
	"errors"
	"net/http"
	"ququiz/lintang/scoring-service/biz/dal/domain"
	"ququiz/lintang/scoring-service/biz/router/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type ScoringService interface {
	GetLeaderboard(ctx context.Context, quizID string) ([]domain.LeaderBoard, error)
}

type LeaderboardHandler struct {
	svc ScoringService
}

func LeaderboardRouter(r *server.Hertz, s ScoringService) {
	handler := &LeaderboardHandler{
		svc: s,
	}

	root := r.Group("/api/v1")
	{
		lH := root.Group("/scoring")
		{
			lH.GET("/:quizID/leaderboard", append(middleware.Protected(), handler.GetLeaderboard)...)
		}
	}
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type getLeaderboardReq struct {
	QuizID string `path:"quizID,required" vd:"regexp('^\\w') && len($) == 24;  msg:'quizID haruslah a-z,A-Z,0-9'"`
}

type getLeaderboardRes struct {
	Leaderboard []domain.LeaderBoard `json:"leaderboard"`
}

func (l *LeaderboardHandler) GetLeaderboard(ctx context.Context, c *app.RequestContext) {
	var req getLeaderboardReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	leaderboard, err := l.svc.GetLeaderboard(ctx, req.QuizID)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, getLeaderboardRes{Leaderboard: leaderboard})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var ierr *domain.Error
	if !errors.As(err, &ierr) {
		return http.StatusInternalServerError
	} else {
		switch ierr.Code() {
		case domain.ErrInternalServerError:
			return http.StatusInternalServerError
		case domain.ErrNotFound:
			return http.StatusNotFound
		case domain.ErrConflict:
			return http.StatusConflict
		case domain.ErrBadParamInput:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}

}
