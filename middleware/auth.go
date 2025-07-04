package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/apps/auth/model"
	"github.com/skyle1995/DevE-Server/database"
	"github.com/skyle1995/DevE-Server/utils/jwt"
	"github.com/spf13/viper"
)

// JWTAuthMiddleware 验证JWT令牌的中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建JWT实例
		jwtConfig := jwt.DefaultConfig()
		// 从配置中获取JWT密钥
		if viper.IsSet("security.jwt_secret") {
			jwtConfig.SigningKey = viper.GetString("security.jwt_secret")
		}
		// 从配置中获取会话超时时间
		if viper.IsSet("server.session_timeout") {
			jwtConfig.ExpiresTime = time.Duration(viper.GetInt("server.session_timeout")) * time.Hour
		}
		jwtInstance := jwt.New(jwtConfig)

		// 尝试从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		tokenString := ""
		var err error

		if authHeader != "" {
			// 从请求头获取令牌
			tokenString, err = jwtInstance.GetTokenFromHeader(authHeader)
			if err != nil {
				// 如果请求头中的令牌无效，尝试从cookie中获取
				tokenCookie, cookieErr := c.Cookie("token")
				if cookieErr != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供认证令牌"})
					c.Abort()
					return
				}
				tokenString, err = jwtInstance.GetTokenFromCookie(tokenCookie)
			}
		} else {
			// 从cookie中获取令牌
			tokenCookie, cookieErr := c.Cookie("token")
			if cookieErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供认证令牌"})
				c.Abort()
				return
			}
			tokenString, err = jwtInstance.GetTokenFromCookie(tokenCookie)
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 解析和验证令牌
		claims, err := jwtInstance.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 检查令牌是否过期
		if jwtInstance.IsTokenExpired(claims) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "认证令牌已过期"})
			c.Abort()
			return
		}

		// 从令牌中获取用户信息
		userID := claims.UserID
		username := claims.Username
		roleStr := claims.Role

		// 将角色字符串转换为整数
		var role int
		if roleStr == "0" {
			role = 0 // 管理员
		} else if roleStr == "2" {
			role = 2 // VIP会员
		} else {
			role = 1 // 普通会员
		}

		// 查询用户是否存在且状态正常
		var user model.User
		result := database.DB.Where("id = ? AND status = ?", userID, 1).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户不存在或已被禁用"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("role", role)

		c.Next()
	}
}

// AdminAuthMiddleware 验证用户是否为管理员的中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取JWTAuthMiddleware中设置的用户角色
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未经授权的访问"})
			c.Abort()
			return
		}

		// 检查用户是否为管理员（角色值为0）
		roleValue, ok := role.(int)
		if !ok || roleValue != 0 {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足，需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
