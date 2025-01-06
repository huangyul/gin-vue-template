package web

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(s *gin.Engine)
}
