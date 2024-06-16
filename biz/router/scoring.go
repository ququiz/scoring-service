package router

import (
	"context"
	"errors"
	"net/http"
	"ququiz/lintang/scoring-service/biz/domain"
	"ququiz/lintang/scoring-service/biz/router/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

type ScoringService interface {
	GetLeaderboard(ctx context.Context, quizID string) ([]domain.LeaderBoard, error)
	RecapQuiz(ctx context.Context, quizID string) error
}

type LeaderboardHandler struct {
	Svc ScoringService
}

func LeaderboardRouter(r *server.Hertz, s ScoringService) {
	handler := &LeaderboardHandler{
		Svc: s,
	}

	root := r.Group("/api/v1")
	{
		lH := root.Group("/scoring")
		{
			lH.GET("/:quizID/leaderboard", append(middleware.Protected(), handler.GetLeaderboard)...)
		}
	}

	dkron := r.Group("/scoring-internal")
	{
		rH := dkron.Group("/recap")
		{
			rH.POST("/:quizID", handler.RecapQuiz)
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
	leaderboard, err := l.Svc.GetLeaderboard(ctx, req.QuizID)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, getLeaderboardRes{Leaderboard: leaderboard})
}

type recapQuizReq struct {
	QuizID string `path:"quizID" vd:"regexp('^\\w') && len($) == 24;  msg:'quizID haruslah a-z,A-Z,0-9'"`
}

type dkronRes struct {
	Message string `json:"message"`
}

func (l *LeaderboardHandler) RecapQuiz(ctx context.Context, c *app.RequestContext) {
	var req recapQuizReq
	err := c.BindAndValidate(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	err = l.Svc.RecapQuiz(ctx, req.QuizID)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dkronRes{Message: "ok"})

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
