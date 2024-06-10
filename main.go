package main

import (
	"context"
	"fmt"
	"net/http"
	"ququiz/lintang/scoring-service/biz/dal/mongodb"
	"ququiz/lintang/scoring-service/biz/dal/rabbitmq"
	"ququiz/lintang/scoring-service/biz/dal/redis"
	"ququiz/lintang/scoring-service/biz/router"
	"ququiz/lintang/scoring-service/biz/service"
	"ququiz/lintang/scoring-service/biz/webapi"
	"ququiz/lintang/scoring-service/config"
	"ququiz/lintang/scoring-service/pkg"
	"time"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/pprof"
	"go.uber.org/zap"
	grpcGo "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		server.WithExitWaitTime(5*time.Second),
	)
	h.Use(accesslog.New(accesslog.WithLogConditionFunc(func(ctx context.Context, c *app.RequestContext) bool {
		if c.FullPath() == "/healthz" {
			return false
		}
		return true
	}))) // jangan pake acess log zap (bikin latency makin tinggi)
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

	// mongo
	mongo := mongodb.NewMongo(cfg)

	quizQueryConsumer := rabbitmq.NewQuizQueryMQConsumer(rmq, leaderboardRedis)
	quizQueryConsumer.ListenAndServe()
	pprof.Register(h)

	// grpc conn
	//grpc
	cc, err := grpcGo.NewClient(cfg.GRPC.AuthGRPCClient+"?wait=30s", grpcGo.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Fatal("Newclient gprc (main)", zap.Error(err))
	}
	defer cc.Close() // close auth grpc connection when graceful shutdown to avoid memory leaks

	// rpc client
	quizQueryClient := webapi.NewQuizQueryClient(cfg)
	authClient := webapi.NewAuthClient(cc)

	// repo
	quizMongoRepo := mongodb.NewQuizRepository(mongo.Conn)

	// producer
	scoringProducer := rabbitmq.NewScoringMQ(rmq)
	quizQueryProducer := rabbitmq.NewQuizQueryProducer(rmq)

	// service
	scoringSvc := service.NewScoringService(leaderboardRedis, quizQueryClient, quizMongoRepo, authClient, scoringProducer, quizQueryProducer)

	// router
	h.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, "service is healthy")
	}) // health probes
	router.LeaderboardRouter(h, scoringSvc)

	// graceful shutdown
	var callback []route.CtxCallback

	callback = append(callback, mongo.Close, rds.Close, rmq.Close) // graceful shutdown
	h.Engine.OnShutdown = append(h.Engine.OnShutdown, callback...)

	h.Spin()
}
