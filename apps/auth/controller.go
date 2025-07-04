package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/middleware"
	"github.com/skyle1995/DevE-Server/utils/captcha"
)

// Controller 处理认证相关的HTTP请求
type Controller struct {
	service *Service
	captcha *captcha.Captcha
}

// NewController 创建一个新的认证控制器实例
func NewController() *Controller {
	return &Controller{
		service: NewService(),
		captcha: captcha.New(nil),
	}
}

// GenerateCaptcha 生成验证码
func (c *Controller) GenerateCaptcha(ctx *gin.Context) {
	// 生成验证码ID
	captchaID := c.captcha.Generate()

	// 将验证码ID存储在会话中
	ctx.SetCookie("captcha_id", captchaID, 600, "/", "", false, true)

	// 输出验证码图片
	if err := c.captcha.WriteImage(ctx.Writer, captchaID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成验证码失败: " + err.Error(),
		})
		return
	}
}

// Login 处理用户登录请求
func (c *Controller) Login(ctx *gin.Context) {
	// 记录请求开始时间
	startTime := time.Now()

	// 绑定请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Captcha  string `json:"captcha"`
		Remember string `json:"remember"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证验证码
	if req.Captcha != "" {
		captchaID, _ := ctx.Cookie("captcha_id")
		if captchaID == "" || !c.captcha.Verify(captchaID, req.Captcha) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "验证码错误",
			})
			return
		}
	}

	// 检查是否记住密码
	remember := false
	if req.Remember == "1" {
		remember = true
	}

	// 调用服务层处理登录
	user, token, err := c.service.Login(req.Username, req.Password, remember)

	// 计算请求处理时间
	latency := time.Since(startTime).Milliseconds()

	if err != nil {
		// 记录登录失败日志
		middleware.LoginLog(0, "登录失败: "+err.Error(), ctx.ClientIP(), ctx.Request.UserAgent(), http.StatusUnauthorized, latency)

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	// 记录登录成功日志
	middleware.LoginLog(user.ID, "用户登录成功", ctx.ClientIP(), ctx.Request.UserAgent(), http.StatusOK, latency)

	// 设置token到cookie
	maxAge := 24 * 3600 // 默认1天
	if remember {
		// 从数据库获取记住密码的有效天数
		var rememberMeDays int = 30 // 默认30天

		// 查询数据库中的设置
		var setting dbmodel.SystemSetting
		result := database.DB.Where("key = ?", "user_remember_me_days").First(&setting)
		if result.Error == nil && setting.Value != "" {
			// 尝试将设置值转换为整数
			fmt.Sscanf(setting.Value, "%d", &rememberMeDays)
		}

		maxAge = rememberMeDays * 24 * 3600 // 设置为数据库配置的天数
	}

	// 设置HttpOnly cookie，增加安全性
	ctx.SetCookie("token", token, maxAge, "/", "", false, true)

	// 返回登录成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

// Register 处理用户注册请求
func (c *Controller) Register(ctx *gin.Context) {
	// 记录请求开始时间
	startTime := time.Now()

	// 绑定请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层处理注册
	err := c.service.Register(req.Username, req.Password, req.Email)

	// 计算请求处理时间
	latency := time.Since(startTime).Milliseconds()

	if err != nil {
		// 记录注册失败日志
		middleware.LoginLog(0, "注册失败: "+err.Error(), ctx.ClientIP(), ctx.Request.UserAgent(), http.StatusBadRequest, latency)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 记录注册成功日志
	middleware.LoginLog(0, "用户注册成功: "+req.Username, ctx.ClientIP(), ctx.Request.UserAgent(), http.StatusOK, latency)

	// 返回注册成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

// GetUserInfo 获取用户信息
func (c *Controller) GetUserInfo(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 调用服务层获取用户信息
	user, err := c.service.GetUserInfo(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	// 返回用户信息
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    user,
	})
}

// UpdatePassword 更新用户密码
func (c *Controller) UpdatePassword(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 绑定请求参数
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层更新密码
	err := c.service.UpdatePassword(userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 返回更新成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码更新成功",
	})
}

// Logout 处理用户注销请求
func (c *Controller) Logout(ctx *gin.Context) {
	// 记录请求开始时间
	startTime := time.Now()

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = uint(0)
	} else if _, ok := userID.(uint); !ok {
		// 如果类型不是uint，则设置为0
		userID = uint(0)
	}

	// 清除token cookie
	ctx.SetCookie("token", "", -1, "/", "", false, true)

	// 计算请求处理时间
	latency := time.Since(startTime).Milliseconds()

	// 记录注销日志
	middleware.LoginLog(userID.(uint), "用户注销成功", ctx.ClientIP(), ctx.Request.UserAgent(), http.StatusOK, latency)

	// 返回注销成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注销成功",
	})
}
