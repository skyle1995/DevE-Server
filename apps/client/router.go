package client

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// SetupClientRoutes 设置客户端路由
func SetupClientRoutes(r *gin.Engine) {
	// 创建控制器实例
	controller := NewController()

	// 客户端API路由组
	clientAPI := r.Group("/api/v1/client")
	clientAPI.Use(middleware.ClientAuthMiddleware())
	{
		// 验证应用
		clientAPI.POST("/verify-app", controller.VerifyApp)

		// 激活卡密
		clientAPI.POST("/activate", controller.ActivateCard)

		// 验证设备
		clientAPI.POST("/verify", controller.VerifyDevice)

		// 换绑卡密
		clientAPI.POST("/rebind", controller.RebindCard)

		// 解绑卡密
		clientAPI.POST("/unbind", controller.UnbindCard)

		// 心跳接口
		clientAPI.POST("/heartbeat", controller.Heartbeat)
	}
}
