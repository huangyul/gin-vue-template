package limiter

import "context"

type Limiter interface {
	limit(ctx context.Context, key string) (bool, error)
}
