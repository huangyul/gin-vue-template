package limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:generate mockgen -source=github.com/redis/go-redis/v9 Cmdable -destination=./mocks/mock_redis.go -package=mocks

//go:embed slide_window.lua
var luaScript string

type RedisSlideWindow struct {
	client   redis.Cmdable
	interval time.Duration
	rate     int
}

func NewRedisSlideWindow(client redis.Cmdable, interval time.Duration, rate int) *RedisSlideWindow {
	return &RedisSlideWindow{client: client, interval: interval, rate: rate}
}

func (r *RedisSlideWindow) Limit(ctx context.Context, key string) (bool, error) {
	return r.client.Eval(ctx, luaScript, []string{key}, r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
