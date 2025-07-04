package notice

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// SetupNoticeRoutes 设置通知相关路由
func SetupNoticeRoutes(router *gin.RouterGroup) {
	// 创建控制器
	controller := NewController()

	// 通知管理路由组
	noticeGroup := router.Group("/notices")

	// 公开接口
	noticeGroup.GET("/dashboard", controller.GetDashboardNotices) // 获取仪表盘公告列表

	// 需要认证的接口
	noticeGroup.Use(middleware.JWTAuthMiddleware())
	{
		// 管理员接口
		noticeGroup.POST("", controller.CreateNotice)                 // 创建通知
		noticeGroup.PUT("/:id", controller.UpdateNotice)              // 更新通知
		noticeGroup.DELETE("/:id", controller.DeleteNotice)           // 删除通知
		noticeGroup.GET("/:id", controller.GetNotice)                 // 获取通知详情
		noticeGroup.GET("", controller.GetNoticeList)                 // 获取通知列表
		noticeGroup.GET("/active", controller.GetActiveNoticeList)    // 获取活动通知列表
		noticeGroup.PUT("/:id/status", controller.UpdateNoticeStatus) // 更新通知状态
	}
}
