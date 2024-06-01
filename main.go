package main

import (
	"fmt"
	"ququiz/lintang/scoring-service/config"
	"ququiz/lintang/scoring-service/pkg"
	"time"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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

	pprof.Register(h)
	// var callback []route.CtxCallback

	// // service & router
	// mongo := mongodb.NewMongo(cfg)
	// rds := rediscache.NewRedis(cfg)

	h.Spin()
}
