package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/response"
)

// AppAuthMiddleware 验证应用API请求的中间件
func AppAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取应用密钥
		appKey := c.GetHeader("X-App-Key")
		appSecret := c.GetHeader("X-App-Secret")

		if appKey == "" || appSecret == "" {
			response.Unauthorized(c, "缺少应用认证信息")
			c.Abort()
			return
		}

		// 查询应用
		var app dbmodel.App
		result := database.DB.Where("app_key = ? AND app_secret = ?", appKey, appSecret).First(&app)
		if result.Error != nil {
			response.Unauthorized(c, "无效的应用认证信息")
			c.Abort()
			return
		}

		// 检查应用状态
		if app.Status != 1 {
			response.Forbidden(c, "应用已被禁用")
			c.Abort()
			return
		}

		// 将应用信息存储到上下文中
		c.Set("app_id", app.ID)
		c.Set("app_name", app.Name)
		c.Set("user_id", app.UserID)

		c.Next()
	}
}

// AppTrialMiddleware 验证应用试用期的中间件
func AppTrialMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取应用ID
		appID, exists := c.Get("app_id")
		if !exists {
			response.Unauthorized(c, "未授权的访问")
			c.Abort()
			return
		}

		// 查询应用
		var app dbmodel.App
		result := database.DB.First(&app, appID)
		if result.Error != nil {
			response.Unauthorized(c, "应用不存在")
			c.Abort()
			return
		}

		// 检查应用是否允许试用
		if app.AllowTrial == 1 && app.TrialAmount > 0 {
			createdTime := app.CreatedAt
			// 根据计费模式处理试用额度
			var trialDuration time.Duration
			if app.BillingMode == 0 { // 时长计费模式
				trialDuration = time.Duration(app.TrialAmount) * time.Hour
			} else { // 点数计费模式
				// 点数计费模式下，假设每个点数对应1小时的使用时间
				trialDuration = time.Duration(app.TrialAmount) * time.Hour
			}
			trialEndTime := createdTime.Add(trialDuration)

			if time.Now().After(trialEndTime) {
				response.Forbidden(c, "应用试用期已过")
				c.Abort()
				return
			}

			// 这里可以添加其他试用相关的逻辑
			// 例如，记录试用使用情况、限制试用功能等
		}

		c.Next()
	}
}
