package ioc

import (
	"github.com/gin-gonic/gin"
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

func InitWebHandler(uHdl *web.UserHandler) []web.Handler {
	return []web.Handler{uHdl}
}

func InitWebMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}
