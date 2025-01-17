package ioc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type config struct {
		Addr     string
		Password string
		DB       int
	}
	var cfg config
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		panic(fmt.Errorf("init redis error: %s", err))
	}
	cmd := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	err := cmd.Ping(context.Background()).Err()
	if err != nil {
		panic(fmt.Errorf("init redis error: %s", err))
	}
	return cmd
}
