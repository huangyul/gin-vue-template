package ioc

import "github.com/gin-gonic/gin"

func InitServer(mdls []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	return server
}

func InitWebMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}
