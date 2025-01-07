package web

import (
	"github.com/gin-gonic/gin"
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

func WriteSuccess(ctx *gin.Context) {
	WriteResponse(ctx, http.StatusOK, 0, "success", nil)
}

func WriteSuccessResponse(ctx *gin.Context, data any) {
	WriteResponse(ctx, http.StatusOK, 0, "", data)
}
