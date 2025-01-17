//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huangyul/gin-vue-template/internal/ioc"
	"github.com/huangyul/gin-vue-template/internal/pkg/ginx/jwt"
	"github.com/huangyul/gin-vue-template/internal/repository"
	"github.com/huangyul/gin-vue-template/internal/repository/dao"
	"github.com/huangyul/gin-vue-template/internal/service"
	"github.com/huangyul/gin-vue-template/internal/web"
)

var UserSet = wire.NewSet(
	jwt.NewHandler,
	dao.NewUserDao,
	repository.NewUserRepository,
	service.NewUserService,
	web.NewUserHandler,
)

var FileSet = wire.NewSet(
	dao.NewFileDao,
	repository.NewFileRepository,
	service.NewFileService,
	web.NewFileHandler,
)

func InitServer() *gin.Engine {
	wire.Build(
		// third party
		ioc.InitDB,
		ioc.InitRedis,

		UserSet,
		FileSet,

		web.NewRouterHandler,

		ioc.InitServer,
		ioc.InitWebHandler,
		ioc.InitWebMiddleware,
	)
	return gin.Default()
}
