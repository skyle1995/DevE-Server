package page

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// RegisterRoutes 注册页面相关路由
func RegisterRoutes(r *gin.Engine) {
	controller := NewController()

	// 需要认证的路由组
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JWTAuthMiddleware())
	{
		// 菜单接口
		apiV1.GET("/menu", controller.GetMenu)
	}
}
