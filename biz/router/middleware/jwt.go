package middleware

import (
	jwt "ququiz/lintang/scoring-service/biz/mw"

	"github.com/cloudwego/hertz/pkg/app"
)

func Protected() []app.HandlerFunc {
	mwJwt := jwt.GetJwtMiddleware()
	mwJwt.MiddlewareInit()
	return []app.HandlerFunc{
		mwJwt.MiddlewareFunc(),
	}

}
