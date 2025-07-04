package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// SetupLogsRoutes 设置日志路由
func SetupLogsRoutes(router *gin.RouterGroup) {
	// 创建控制器实例
	controller := NewController()

	// 日志管理路由组
	logsGroup := router.Group("/logs")
	{
		// 获取日志列表 - 管理员可查看所有日志
		logsGroup.GET("", middleware.AdminAuthMiddleware(), controller.GetLogs)
		// 清空日志 - 仅管理员可操作
		logsGroup.DELETE("", middleware.AdminAuthMiddleware(), controller.ClearLogs)
		// 获取用户自己的日志 - 普通用户可查看自己的日志
		logsGroup.GET("/my", middleware.JWTAuthMiddleware(), controller.GetMyLogs)
	}
}
