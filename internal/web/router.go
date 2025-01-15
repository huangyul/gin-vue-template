package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Meta 定义路由元信息结构
type Meta struct {
	Title string   `json:"title"`
	Icon  string   `json:"icon,omitempty"`
	Rank  int      `json:"rank,omitempty"`
	Roles []string `json:"roles,omitempty"`
	Auths []string `json:"auths,omitempty"`
}

// Route 定义路由结构
type Route struct {
	Path      string  `json:"path"`
	Name      string  `json:"name,omitempty"`
	Component string  `json:"component,omitempty"`
	Meta      Meta    `json:"meta"`
	Children  []Route `json:"children,omitempty"`
}

// 定义 PermissionRouter
var permissionRouter = Route{
	Path: "/permission",
	Meta: Meta{
		Title: "权限管理",
		Icon:  "ep:lollipop",
		Rank:  10,
	},
	Children: []Route{
		{
			Path: "/permission/page/index",
			Name: "PermissionPage",
			Meta: Meta{
				Title: "页面权限",
				Roles: []string{"admin", "common"},
			},
		},
		{
			Path: "/permission/button",
			Meta: Meta{
				Title: "按钮权限",
				Roles: []string{"admin", "common"},
			},
			Children: []Route{
				{
					Path:      "/permission/button/router",
					Component: "permission/button/index",
					Name:      "PermissionButtonRouter",
					Meta: Meta{
						Title: "路由返回按钮权限",
						Auths: []string{
							"permission:btn:add",
							"permission:btn:edit",
							"permission:btn:delete",
						},
					},
				},
				{
					Path:      "/permission/button/login",
					Component: "permission/button/perms",
					Name:      "PermissionButtonLogin",
					Meta: Meta{
						Title: "登录接口返回按钮权限",
					},
				},
			},
		},
	},
}

var userRouter = Route{
	Path: "/user",
	Meta: Meta{
		Title: "用户管理",
		Icon:  "ep:user",
		Rank:  20,
	},
	Children: []Route{
		{
			Path: "/user/page/list",
			Name: "UserPageList",
			Meta: Meta{
				Title: "用户管理",
				Roles: []string{"admin", "common"},
			},
		},
	},
}

var fileRouter = Route{
	Path: "/file",
	Meta: Meta{
		Title: "文件管理",
		Icon:  "ep:files",
		Rank:  20,
	},
	Children: []Route{
		{
			Path: "/file/page/list",
			Name: "FilePageList",
			Meta: Meta{
				Title: "文件管理",
				Roles: []string{"admin", "common"},
			},
		},
	},
}

type RouterHandler struct{}

func NewRouterHandler() *RouterHandler {
	return &RouterHandler{}
}

func (r *RouterHandler) RegisterRoutes(s *gin.Engine) {
	s.GET("/get-async-routes", r.GetRoutes)
}

func (r *RouterHandler) GetRoutes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    []Route{permissionRouter, userRouter, fileRouter},
	})
}
