package middleware

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/crypto"
	"github.com/skyle1995/DevE-Server/utils/response"
	"github.com/skyle1995/DevE-Server/utils/timeutil"
	"github.com/spf13/viper"
)

// ClientAuthMiddleware 客户端认证中间件
// 用于验证客户端请求的合法性
// 客户端请求需要在Header中携带以下信息：
// - App-Key: 应用的AppKey
// - App-Sign: 请求签名，使用AppSecret对请求参数进行签名
// - Timestamp: 请求时间戳，用于防止重放攻击
func ClientAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的认证信息
		appKey := c.GetHeader("App-Key")
		appSign := c.GetHeader("App-Sign")
		timestamp := c.GetHeader("Timestamp")

		// 检查必要的认证信息是否存在
		if appKey == "" || appSign == "" || timestamp == "" {
			response.FailWithDetailed(gin.H{
				"reload": true,
			}, "缺少必要的认证信息", c)
			c.Abort()
			return
		}

		// 验证时间戳
		if err := validateTimestamp(timestamp); err != nil {
			response.FailWithDetailed(gin.H{
				"reload": true,
			}, err.Error(), c)
			c.Abort()
			return
		}

		// 查询应用信息
		var app dbmodel.App
		result := database.DB.Where("app_key = ?", appKey).First(&app)
		if result.Error != nil {
			response.FailWithDetailed(gin.H{
				"reload": true,
			}, "应用不存在或已被禁用", c)
			c.Abort()
			return
		}

		// 检查应用状态
		if app.Status != 1 {
			response.FailWithDetailed(gin.H{
				"reload": true,
			}, "应用已被禁用", c)
			c.Abort()
			return
		}

		// 验证请求签名
		err := ValidateAppSign(c, app)
		if err != nil {
			response.FailWithDetailed(gin.H{
				"reload": true,
			}, "签名验证失败: "+err.Error(), c)
			c.Abort()
			return
		}

		// 将应用信息存储到上下文中，供后续处理使用
		c.Set("app", app)
		c.Next()
	}
}

// ValidateAppSign 验证应用签名
// 签名规则：
// 1. 将请求参数按照参数名的字典序排序
// 2. 将参数名和参数值拼接成字符串，格式为：参数名=参数值
// 3. 将拼接后的字符串用&连接，再加上应用的AppSecret
// 4. 对最终的字符串进行MD5加密，得到签名
func ValidateAppSign(c *gin.Context, app dbmodel.App) error {
	// 获取请求参数
	params := make(map[string]string)

	// 获取URL查询参数
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	// 获取表单参数
	c.Request.ParseForm()
	for k, v := range c.Request.PostForm {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	// 获取请求头中的签名
	sign := c.GetHeader("App-Sign")
	if sign == "" {
		return errors.New("缺少签名信息")
	}

	// 使用utils/crypto包中的函数验证签名
	if !crypto.VerifySign(params, sign, app.AppSecret) {
		return errors.New("签名验证失败")
	}

	return nil
}

// validateTimestamp 验证时间戳
// timestamp: 请求时间戳（秒级Unix时间戳）
func validateTimestamp(timestamp string) error {
	// 将时间戳转换为整数
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return errors.New("时间戳格式错误")
	}

	// 获取当前时间
	now := timeutil.Now()

	// 将时间戳转换为时间对象
	timestampTime := time.Unix(timestampInt, 0)

	// 获取配置的时间戳有效期（秒）
	timestampExpire := viper.GetInt64("security.timestamp_expire")
	if timestampExpire <= 0 {
		// 默认5分钟
		timestampExpire = 300
	}

	// 计算时间差
	timeDiff := timeutil.DiffSeconds(now, timestampTime)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}

	// 检查时间戳是否在有效期内
	if timeDiff > timestampExpire {
		return errors.New("时间戳已过期")
	}

	return nil
}
