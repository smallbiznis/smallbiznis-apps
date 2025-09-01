package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/smallbiznis/smallbiznis-apps/pkg/config"
	"go.uber.org/fx"
)

var Module = fx.Module("redis",
	fx.Provide(New),
)

func New(lc fx.Lifecycle, c *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:        c.Redis.Addr,
		Password:    c.Redis.Password, // no password set
		DB:          c.Redis.DB,       // use default DB
		PoolSize:    c.Redis.PoolSize,
		PoolTimeout: c.Redis.PoolTimeout,
	})

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return rdb.Close()
		},
	})

	return rdb
}
