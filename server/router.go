package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apps "github.com/skyle1995/DevE-Server/apps/app"
	"github.com/skyle1995/DevE-Server/apps/auth"
	"github.com/skyle1995/DevE-Server/apps/card"
	"github.com/skyle1995/DevE-Server/apps/client"
	"github.com/skyle1995/DevE-Server/apps/logs"
	"github.com/skyle1995/DevE-Server/apps/notice"
	"github.com/skyle1995/DevE-Server/apps/page"
	"github.com/skyle1995/DevE-Server/apps/setting"
	"github.com/skyle1995/DevE-Server/apps/user"
	"github.com/skyle1995/DevE-Server/middleware"
	"github.com/skyle1995/DevE-Server/public"
	"github.com/spf13/viper"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 设置运行模式
	mode := viper.GetString("server.mode")
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 创建路由引擎
	r := gin.New()

	// 使用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 映射前端静态资源
	// 将dist/index.html映射到/
	r.GET("/", func(c *gin.Context) {
		indexHTML, err := public.Public.ReadFile("dist/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "无法读取index.html")
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(indexHTML))
	})

	// 将dist/favicon.ico映射到/favicon.ico
	r.GET("/favicon.ico", func(c *gin.Context) {
		favicon, err := public.Public.ReadFile("dist/favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "image/x-icon", favicon)
	})

	// 将dist/assets映射到/assets
	r.GET("/assets/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		filepath = strings.TrimPrefix(filepath, "/")

		file, err := public.Public.ReadFile("dist/assets/" + filepath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		// 根据文件扩展名设置Content-Type
		contentType := "application/octet-stream"
		if strings.HasSuffix(filepath, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(filepath, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(filepath, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(filepath, ".jpg") || strings.HasSuffix(filepath, ".jpeg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(filepath, ".svg") {
			contentType = "image/svg+xml"
		} else if strings.HasSuffix(filepath, ".json") {
			contentType = "application/json"
		}

		c.Data(http.StatusOK, contentType, file)
	})

	// 创建控制器实例
	authController := auth.NewController()

	// 公开路由组
	public := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := public.Group("/auth")
		{
			// 登录路由添加登录日志中间件
			auth.POST("/login", middleware.LoginLogMiddleware(), authController.Login)
			auth.POST("/register", authController.Register)
			auth.GET("/captcha", authController.GenerateCaptcha)
			auth.POST("/logout", authController.Logout)
		}
	}

	// 设置系统设置路由
	setting.SetupSettingRoutes(r)

	// 设置通知路由
	notice.SetupNoticeRoutes(r.Group("/api/v1"))

	// 设置用户路由
	user.SetupUserRoutes(r)

	// 设置卡密路由
	card.SetupCardRoutes(r)

	// 设置客户端路由
	client.SetupClientRoutes(r)

	// 设置应用路由
	apps.SetupAppsRoutes(r)

	// 设置应用API路由
	apps.SetupAppAPIRoutes(r)

	// 设置日志路由
	logs.SetupLogsRoutes(r.Group("/api/v1"))

	// 设置页面路由
	page.RegisterRoutes(r)

	// 需要认证的路由组
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// 用户相关路由
		user := protected.Group("/user")
		{
			user.GET("/info", authController.GetUserInfo)
			user.POST("/password", authController.UpdatePassword)
		}

		// 页面相关路由已移至page模块的router.go文件中

		// 管理员路由组
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware())
		{
			// 这里可以添加管理员特有的路由
		}
	}

	return r
}
