package user

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// SetupUserRoutes 设置用户相关路由
func SetupUserRoutes(r *gin.Engine) {
	// 创建用户控制器实例
	userController := NewController()

	// 需要认证的路由组
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// 用户个人资料相关路由（所有已登录用户可访问）
		profile := protected.Group("/profile")
		{
			profile.GET("", userController.GetUserProfile)    // 获取个人资料
			profile.PUT("", userController.UpdateUserProfile) // 更新个人资料
		}

		// 用户权限相关路由（所有已登录用户可访问）
		protected.GET("/permission", userController.GetUserPermission) // 获取当前用户权限

		// 管理员路由组（仅管理员可访问）
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware())
		{
			// 用户管理相关路由
			users := admin.Group("/users")
			{
				users.GET("", userController.GetUserList)       // 获取用户列表
				users.GET("/:id", userController.GetUserByID)   // 获取用户详情
				users.POST("", userController.CreateUser)       // 创建用户
				users.PUT("/:id", userController.UpdateUser)    // 更新用户
				users.DELETE("/:id", userController.DeleteUser) // 删除用户
			}
		}
	}
}
