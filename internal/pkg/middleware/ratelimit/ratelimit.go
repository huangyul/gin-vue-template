package ratelimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/pkg/errno"
	"github.com/huangyul/gin-vue-template/internal/pkg/limiter"
	"net/http"
)

type Builder struct {
	l      limiter.Limiter
	prefix string
}

func NewBuilder(l limiter.Limiter) *Builder {
	return &Builder{l: l, prefix: "ip-limiter"}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		limited, err := b.limit(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errno.InternalServerError.SetMessage(err.Error()))
			return
		}
		if limited {
			c.AbortWithError(http.StatusTooManyRequests, errno.BadRequest.SetMessage("to many request"))
			return
		}
		c.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s_%s", b.prefix, ctx.ClientIP())
	return b.l.Limit(ctx, key)
}
