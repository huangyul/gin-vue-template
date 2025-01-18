package limiter

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRedisSlideWindow_Limit(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "",
		DB:       1,
	})
	err := client.Ping(context.Background()).Err()
	assert.Nil(t, err)
	err = client.FlushDB(context.Background()).Err()
	assert.NoError(t, err)
	interval := time.Minute
	rate := 20
	key := "test_key"
	limiter := NewRedisSlideWindow(client, interval, rate)

	res, err := limiter.Limit(context.Background(), key)
	assert.Equal(t, redis.Nil, err)
	assert.False(t, res)
	for i := 0; i < 30; i++ {
		limiter.Limit(context.Background(), key)
	}
	res, err = limiter.Limit(context.Background(), key)
	assert.NoError(t, err)
	assert.True(t, res)
	err = client.FlushDB(context.Background()).Err()
	assert.NoError(t, err)
}
