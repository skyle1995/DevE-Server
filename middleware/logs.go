package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
)

// LoginLogMiddleware 登录日志中间件
func LoginLogMiddleware() gin.HandlerFunc {
	return LogMiddleware(dbmodel.LogTypeLogin)
}

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return LogMiddleware(dbmodel.LogTypeOperation)
}

// SystemLogMiddleware 系统日志中间件
func SystemLogMiddleware() gin.HandlerFunc {
	return LogMiddleware(dbmodel.LogTypeSystem)
}

// LogMiddleware 通用日志中间件
func LogMiddleware(logType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重置请求体，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建响应体记录器
		responseWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		// 处理请求
		c.Next()

		// 计算请求处理时间
		latency := time.Since(startTime)

		// 创建日志记录
		log := dbmodel.Logs{
			Type:         logType,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			StatusCode:   c.Writer.Status(),
			Latency:      latency.Milliseconds(),
			RequestBody:  string(requestBody),
			ResponseBody: responseWriter.body.String(),
		}

		// 根据日志类型设置相关ID
		switch logType {
		case dbmodel.LogTypeApp:
			// 获取应用ID
			appID, exists := c.Get("app_id")
			if !exists {
				return // 如果没有应用ID，不记录日志
			}
			log.AppID = appID.(uint)
		case dbmodel.LogTypeLogin:
			// 登录日志不需要检查用户ID，因为登录时用户ID可能还不存在
			// 用户ID将在登录成功后由 LoginLog 函数设置
		case dbmodel.LogTypeOperation:
			// 获取用户ID
			userID, exists := c.Get("user_id")
			if !exists {
				return // 如果没有用户ID，不记录日志
			}
			log.UserID = userID.(uint)
		}

		// 异步保存日志
		go func(log dbmodel.Logs) {
			database.DB.Create(&log)
		}(log)
	}
}

// responseBodyWriter 响应体记录器
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法，记录响应体
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// SystemLog 记录系统日志
func SystemLog(content string) {
	log := dbmodel.Logs{
		Type:    dbmodel.LogTypeSystem,
		Content: content,
	}
	database.DB.Create(&log)
}

// LoginLog 记录登录日志
func LoginLog(userID uint, content string, ip string, userAgent string, statusCode int, latency ...int64) {
	// 创建日志记录
	log := dbmodel.Logs{
		Type:       dbmodel.LogTypeLogin,
		UserID:     userID,
		Content:    content,
		IP:         ip,
		UserAgent:  userAgent,
		StatusCode: statusCode,
		Latency:    0, // 默认值
	}

	// 如果提供了耗时参数，则使用提供的值
	if len(latency) > 0 && latency[0] > 0 {
		log.Latency = latency[0]
	}

	database.DB.Create(&log)
}
