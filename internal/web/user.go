package web

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/service"
	"net/http"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) Register(server *gin.Engine) {
	ug := server.Group("/user")
	ug.GET("/login", h.Login)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	ctx.String(http.StatusOK, ctx.Query("name"))
}
