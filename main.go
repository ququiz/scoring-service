package main

import (
	"fmt"
	"ququiz/lintang/scoring-service/biz/dal/rabbitmq"
	"ququiz/lintang/scoring-service/biz/dal/redis"
	"ququiz/lintang/scoring-service/biz/router"
	"ququiz/lintang/scoring-service/biz/service"
	"ququiz/lintang/scoring-service/config"
	"ququiz/lintang/scoring-service/pkg"
	"time"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/pprof"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		hlog.Fatalf("Config error: %s", err)
	}
	gopool.SetCap(500000) // naikin goroutine netpool for high performance

	logsCores := pkg.InitZapLogger(cfg)
	defer logsCores.Sync()
	hlog.SetLogger(logsCores)

	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(`0.0.0.0:%s`, cfg.HTTP.Port)),
		server.WithExitWaitTime(4*time.Second),
	)
	h.Use(accesslog.New()) // jangan pake acess log zap (bikin latency makin tinggi)
	// setup cors
	corsHandler := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "content-type", "authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"},
		ExposeHeaders:    []string{"Origin", "content-type", "authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"},
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	})

	h.Use(corsHandler)

	// redis
	rds := redis.NewRedis(cfg)
	leaderboardRedis := redis.NewLeaderboardRedis(rds.Client)

	// rabbitmq
	rmq := rabbitmq.NewRabbitMQ(cfg)

	quizQueryConsumer := rabbitmq.NewQuizQueryMQConsumer(rmq, leaderboardRedis)
	quizQueryConsumer.ListenAndServe()
	pprof.Register(h)

	// service
	scoringSvc := service.NewScoringService(leaderboardRedis)

	// router
	router.LeaderboardRouter(h, scoringSvc)

	// var callback []route.CtxCallback

	// // service & router
	// mongo := mongodb.NewMongo(cfg)
	// rds := rediscache.NewRedis(cfg)

	h.Spin()
}
