package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/pkg/ginx/jwt"
	"github.com/huangyul/gin-vue-template/internal/pkg/middleware"
	"github.com/huangyul/gin-vue-template/internal/web"
)

func InitServer(mdls []gin.HandlerFunc, hdls []web.Handler) *gin.Engine {
	server := gin.Default()
	server.Static("/static", "./static")
	server.Use(mdls...)
	for _, hdl := range hdls {
		hdl.RegisterRoutes(server)
	}
	return server
}

func InitWebHandler(uHdl *web.UserHandler, routerHdl *web.RouterHandler, fHdl *web.FileHandler) []web.Handler {
	return []web.Handler{uHdl, routerHdl, fHdl}
}

func InitWebMiddleware(jwtHdl *jwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// CORS
		cors.New(cors.Config{
			AllowMethods:  []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:  []string{"*"},
			AllowOrigins:  []string{"*"},
			ExposeHeaders: []string{"Content-Length"},
			MaxAge:        86400,
		}),
		// login
		middleware.NewJWTMiddlewareBuild(jwtHdl).AddWhiteList("/user/login", "/user/refresh-token", "/user/register").Build(),
	}
}
