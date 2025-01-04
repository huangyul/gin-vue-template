//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huangyul/gin-vue-template/internal/ioc"
)

func InitServer() *gin.Engine {
	wire.Build(
		ioc.InitServer,
		ioc.InitWebMiddleware)
	return gin.Default()
}
