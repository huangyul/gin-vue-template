package web

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/pkg/errno"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteResponse(ctx *gin.Context, httpCode int, code int, message string, data interface{}) {
	ctx.JSON(httpCode, APIResponse{Code: code, Message: message, Data: data})
}

func WriteError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, APIResponse{Code: code, Message: message})
}

func WriteErrno(ctx *gin.Context, err *errno.Errno) {
	WriteError(ctx, err.Code, err.Message)
}
func WriteSuccess(ctx *gin.Context) {
	WriteResponse(ctx, http.StatusOK, 0, "success", nil)
}

func WriteSuccessResponse(ctx *gin.Context, data any) {
	WriteResponse(ctx, http.StatusOK, 0, "", data)
}
