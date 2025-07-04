package apps

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// SetupAppsRoutes 设置应用相关路由
func SetupAppsRoutes(r *gin.Engine) {
	// 创建控制器实例
	controller := NewController()

	// 需要认证的路由组
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// 应用相关路由（所有已登录用户可访问）
		apps := protected.Group("/apps")
		{
			apps.POST("", controller.CreateApp)                      // 创建应用
			apps.GET("", controller.GetAppList)                      // 获取应用列表
			apps.GET("/:id", controller.GetAppByID)                  // 获取应用详情
			apps.PUT("/:id", controller.UpdateApp)                   // 更新应用
			apps.DELETE("/:id", controller.DeleteApp)                // 删除应用
			apps.POST("/:id/secret", controller.RegenerateAppSecret) // 重新生成应用密钥
		}

	}
}
