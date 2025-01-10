package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/gin-vue-template/internal/dto"
	"github.com/huangyul/gin-vue-template/internal/pkg/ginx/jwt"
	"github.com/huangyul/gin-vue-template/internal/service"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
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
	ug.POST("/list", h.List)
	ug.GET("/detail/:id", h.Detail)
	ug.GET("/delete/:id", h.Delete)
	ug.POST("/update", h.Update)
	ug.POST("/create", h.Create)
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

func (h *UserHandler) List(ctx *gin.Context) {
	var r dto.UserListQueryParams
	if err := ctx.ShouldBind(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	us, count, err := h.svc.List(ctx, r)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var users []dto.UserResp
	for _, u := range us {
		users = append(users, dto.UserResp{
			Username:  u.Username,
			Nickname:  u.Nickname,
			ID:        u.ID,
			CreatedAt: u.CreatedAt.Format(time.DateOnly),
		})
	}
	WriteSuccessResponse(ctx, dto.UserListResp{
		Data:  users,
		Total: count,
	})

}

func (h *UserHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "非法id"})
		return
	}
	u, err := h.svc.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	WriteSuccessResponse(ctx, dto.UserResp{
		Username:  u.Username,
		Nickname:  u.Nickname,
		ID:        u.ID,
		CreatedAt: u.CreatedAt.Format(time.DateOnly),
	})
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.svc.DeleteByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	WriteSuccess(ctx)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	type req struct {
		ID       int64  `json:"id" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
	}
	var r req
	if err := ctx.ShouldBind(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.svc.Update(ctx, r.ID, r.Nickname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	WriteSuccess(ctx)
}

func (h *UserHandler) Create(ctx *gin.Context) {
	type req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
	}
	var r req
	if err := ctx.ShouldBind(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.svc.Create(ctx, r.Username, r.Password, r.Nickname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	WriteSuccess(ctx)
}
