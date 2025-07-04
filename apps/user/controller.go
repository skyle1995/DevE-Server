package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller 处理用户相关的HTTP请求
type Controller struct {
	service *Service
}

// NewController 创建一个新的用户控制器实例
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// GetUserList 获取用户列表（仅管理员可用）
func (c *Controller) GetUserList(ctx *gin.Context) {
	// 获取查询参数
	role := ctx.DefaultQuery("role", "")
	status := ctx.DefaultQuery("status", "")
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "20")

	// 调用服务获取用户列表
	users, total, err := c.service.GetUserList(role, status, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户列表失败: " + err.Error(),
		})
		return
	}

	// 返回用户列表
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"total": total,
			"items": users,
		},
	})
}

// GetUserByID 根据ID获取用户信息（仅管理员可用）
func (c *Controller) GetUserByID(ctx *gin.Context) {
	// 获取用户ID
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少用户ID",
		})
		return
	}

	// 调用服务获取用户信息
	user, err := c.service.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户信息失败: " + err.Error(),
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

// CreateUser 创建新用户（仅管理员可用）
func (c *Controller) CreateUser(ctx *gin.Context) {
	// 绑定请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Role     string `json:"role" binding:"required"`
		Status   string `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务创建用户
	err := c.service.CreateUser(req.Username, req.Password, req.Email, req.Role, req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "创建用户失败: " + err.Error(),
		})
		return
	}

	// 返回创建成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

// UpdateUser 更新用户信息（仅管理员可用）
func (c *Controller) UpdateUser(ctx *gin.Context) {
	// 获取用户ID
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少用户ID",
		})
		return
	}

	// 绑定请求参数
	var req struct {
		Email  string `json:"email" binding:"omitempty,email"`
		Role   string `json:"role"`
		Status string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务更新用户
	err := c.service.UpdateUser(userID, req.Email, req.Role, req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "更新用户失败: " + err.Error(),
		})
		return
	}

	// 返回更新成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteUser 删除用户（仅管理员可用）
func (c *Controller) DeleteUser(ctx *gin.Context) {
	// 获取用户ID
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少用户ID",
		})
		return
	}

	// 调用服务删除用户
	err := c.service.DeleteUser(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "删除用户失败: " + err.Error(),
		})
		return
	}

	// 返回删除成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// GetUserProfile 获取当前用户的个人资料
func (c *Controller) GetUserProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 调用服务获取用户资料
	profile, err := c.service.GetUserProfile(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取个人资料失败: " + err.Error(),
		})
		return
	}

	// 返回用户资料
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    profile,
	})
}

// UpdateUserProfile 更新当前用户的个人资料
func (c *Controller) UpdateUserProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 绑定请求参数
	var req struct {
		Email string `json:"email" binding:"omitempty,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务更新用户资料
	err := c.service.UpdateUserProfile(userID.(uint), req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "更新个人资料失败: " + err.Error(),
		})
		return
	}

	// 返回更新成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// GetUserPermission 获取当前用户权限
func (c *Controller) GetUserPermission(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权的访问"})
		return
	}

	// 获取用户权限
	permission, err := c.service.GetUserPermission(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取权限失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    permission,
	})
}
