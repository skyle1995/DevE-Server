package setting

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// SetupSettingRoutes 设置系统设置相关路由
func SetupSettingRoutes(r *gin.Engine) {
	// 创建系统设置控制器实例
	settingController := NewController()

	// 公开路由（无需认证）
	public := r.Group("/api/v1")
	{
		// 站点信息路由
		public.GET("/site/info", settingController.GetSiteInfo) // 获取站点信息
	}

	// 需要认证的路由组
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// 系统设置相关路由（所有已登录用户可访问）
		settings := protected.Group("/settings")
		{
			settings.GET("", settingController.GetSettings) // 获取系统设置列表
		}

		// 管理员路由组（仅管理员可访问）
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware())
		{
			// 系统设置管理相关路由
			adminSettings := admin.Group("/settings")
			{
				adminSettings.GET("/:key", settingController.GetSettingByKey)  // 根据键名获取设置
				adminSettings.PUT("", settingController.UpdateSetting)         // 更新设置
				adminSettings.POST("", settingController.CreateSetting)        // 创建设置
				adminSettings.DELETE("/:key", settingController.DeleteSetting) // 删除设置
			}
		}
	}
}
