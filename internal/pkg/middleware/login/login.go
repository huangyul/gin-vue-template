package login

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/pkg/ginx/jwt"
	"net/http"
	"strings"
)

type JWTMiddlewareBuild struct {
	WhiteList []string
	jwtHdl    *jwt.Handler
}

func NewJWTMiddlewareBuild(jwtHdl *jwt.Handler) *JWTMiddlewareBuild {
	return &JWTMiddlewareBuild{jwtHdl: jwtHdl}
}

func (l *JWTMiddlewareBuild) AddWhiteList(whiteList ...string) *JWTMiddlewareBuild {
	l.WhiteList = append(l.WhiteList, whiteList...)
	return l
}

func (l *JWTMiddlewareBuild) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, p := range l.WhiteList {
			if p == ctx.Request.URL.Path {
				ctx.Next()
			}
		}
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenString, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claim, err := l.jwtHdl.ParseToken(segs[1])
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("userId", claim.UserId)
		ctx.Next()
	}
}
