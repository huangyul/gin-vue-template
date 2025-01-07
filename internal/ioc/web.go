package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/pkg/ginx/jwt"
	"github.com/huangyul/gin-vue-template/internal/web"
)

func InitServer(mdls []gin.HandlerFunc, hdls []web.Handler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	for _, hdl := range hdls {
		hdl.RegisterRoutes(server)
	}
	return server
}

func InitWebHandler(uHdl *web.UserHandler, routerHdl *web.RouterHandler) []web.Handler {
	return []web.Handler{uHdl, routerHdl}
}

func InitWebMiddleware(jwtHdl *jwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// login
		web.NewLoginMiddlewareBuild(jwtHdl).AddWhiteList("/user/login", "/user/refresh-token").Build(),
	}
}
