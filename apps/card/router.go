package card

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/middleware"
)

// SetupCardRoutes 设置卡密相关路由
func SetupCardRoutes(r *gin.Engine) {
	cardController := NewController()

	// 用户API路由组
	cardGroup := r.Group("/api/v1/card")
	cardGroup.Use(middleware.JWTAuthMiddleware())
	{
		// 卡密类型管理
		cardGroup.GET("/types", cardController.GetCardTypes)          // 获取卡密类型列表
		cardGroup.POST("/types", cardController.CreateCardType)       // 创建卡密类型
		cardGroup.PUT("/types/:id", cardController.UpdateCardType)    // 更新卡密类型
		cardGroup.DELETE("/types/:id", cardController.DeleteCardType) // 删除卡密类型

		// 卡密管理
		cardGroup.GET("/cards", cardController.GetCards)          // 获取卡密列表
		cardGroup.POST("/generate", cardController.GenerateCards) // 生成卡密
		cardGroup.PUT("/cards/:id", cardController.UpdateCard)    // 更新卡密
		cardGroup.DELETE("/cards/:id", cardController.DeleteCard) // 删除卡密
	}
}
