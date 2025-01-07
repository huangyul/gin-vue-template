package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/pkg/ginx/jwt"
	"github.com/huangyul/gin-vue-template/internal/service"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

var (
	emailPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordPattern = `^[a-zA-Z0-9]{6,18}$`
)

type UserHandler struct {
	svc            service.UserService
	jwt            *jwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, jwt *jwt.Handler) *UserHandler {
	return &UserHandler{
		svc:            svc,
		jwt:            jwt,
		emailRexExp:    regexp.MustCompile(emailPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordPattern, regexp.None),
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/login", h.Login)
	ug.POST("/register", h.Register)
	ug.POST("/refresh-token", h.RefreshToken)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var egg errgroup.Group
	var token string
	var refreshToken string
	egg.Go(func() error {
		var er error
		token, er = h.jwt.GenerateToken(user.ID)
		return er
	})
	egg.Go(func() error {
		var er error
		refreshToken, er = h.jwt.GenerateRefreshToken(user.ID)
		return er
	})
	if err := egg.Wait(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type LoginResp struct {
		Avatar       string   `json:"avatar"`       // 头像
		Username     string   `json:"username"`     // 用户名
		Nickname     string   `json:"nickname"`     // 昵称
		Roles        []string `json:"roles"`        // 当前登录用户的角色
		Permissions  []string `json:"permissions"`  // 按钮级别权限
		AccessToken  string   `json:"accessToken"`  // `token`
		RefreshToken string   `json:"refreshToken"` // 用于调用刷新`accessToken`的接口时所需的`token`
		Expires      string   `json:"expires"`      // `accessToken`的过期时间（格式'xxxx/xx/xx xx:xx:xx'）
	}
	WriteSuccessResponse(ctx, LoginResp{
		Roles:        []string{"admin"},
		Permissions:  []string{"*:*:*"},
		Nickname:     "小红",
		Avatar:       "https://avatars.githubusercontent.com/u/44761321",
		Username:     user.Username,
		AccessToken:  token,
		RefreshToken: refreshToken,
		Expires:      time.Now().Add(time.Hour * 24 * 2).Format("2006/01/02 15:04:05"),
	})
}

func (h *UserHandler) Register(ctx *gin.Context) {
	type RegisterRequest struct {
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}
	var req RegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password must be between 6 and 18 digits"})
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "passwords do not match"})
		return
	}
	err = h.svc.Register(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": req.Username})
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	type req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	var r req
	if err := ctx.ShouldBind(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.jwt.Refresh(r.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type RefreshTokenResp struct {
		AccessToken string `json:"accessToken"`
	}
	WriteSuccessResponse(ctx, RefreshTokenResp{
		AccessToken: token,
	})
}
