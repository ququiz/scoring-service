package rediscache

import (
	"context"
	"ququiz/lintang/scoring-service/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(cfg *config.Config) *Redis {
	// ini kalo deploy ke kubernetes
	// redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{"redis1.redis-svc.redis.svc.cluster.local:16379", ""},
	// })
	cli := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.RedisAddr,
		// Password: cfg.Redis.RedisPassword,
		DB: 0,
	})

	return &Redis{cli}
}

func (r *Redis) Close(ctx context.Context) {
	r.Client.Close()
}
