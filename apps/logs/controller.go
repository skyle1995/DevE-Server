package logs

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/response"
)

// Controller 日志控制器
type Controller struct{}

// NewController 创建日志控制器实例
func NewController() *Controller {
	return &Controller{}
}

// GetLogs 获取日志列表
func (c *Controller) GetLogs(ctx *gin.Context) {
	// 获取日志类型
	logType, err := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	if err != nil {
		response.BadRequest(ctx, "无效的日志类型")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// 查询日志
	var logs []dbmodel.Logs
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询条件
	query := database.DB.Model(&dbmodel.Logs{})

	// 根据日志类型筛选
	if logType > 0 {
		query = query.Where("type = ?", logType)
	}

	// 应用ID筛选
	if appID := ctx.Query("app_id"); appID != "" {
		appIDInt, err := strconv.ParseUint(appID, 10, 32)
		if err == nil {
			query = query.Where("app_id = ?", uint(appIDInt))
		}
	}

	// 用户ID筛选
	if userID := ctx.Query("user_id"); userID != "" {
		userIDInt, err := strconv.ParseUint(userID, 10, 32)
		if err == nil {
			query = query.Where("user_id = ?", uint(userIDInt))
		}
	}

	// 查询总数
	result := query.Count(&total)
	if result.Error != nil {
		response.InternalServerError(ctx, "获取日志失败: "+result.Error.Error())
		return
	}

	// 分页查询
	result = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)
	if result.Error != nil {
		response.InternalServerError(ctx, "获取日志失败: "+result.Error.Error())
		return
	}

	// 返回日志列表
	response.SuccessWithPagination(ctx, "获取日志成功", logs, page, pageSize, total)
}

// GetMyLogs 获取当前用户的日志列表
func (c *Controller) GetMyLogs(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	// 获取日志类型
	logType, err := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	if err != nil {
		response.BadRequest(ctx, "无效的日志类型")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// 查询日志
	var logs []dbmodel.Logs
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询条件 - 只查询当前用户的日志
	query := database.DB.Model(&dbmodel.Logs{}).Where("user_id = ?", userID)

	// 根据日志类型筛选
	if logType > 0 {
		query = query.Where("type = ?", logType)
	}

	// 应用ID筛选
	if appID := ctx.Query("app_id"); appID != "" {
		appIDInt, err := strconv.ParseUint(appID, 10, 32)
		if err == nil {
			query = query.Where("app_id = ?", uint(appIDInt))
		}
	}

	// 查询总数
	result := query.Count(&total)
	if result.Error != nil {
		response.InternalServerError(ctx, "获取日志失败: "+result.Error.Error())
		return
	}

	// 分页查询
	result = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)
	if result.Error != nil {
		response.InternalServerError(ctx, "获取日志失败: "+result.Error.Error())
		return
	}

	// 返回日志列表
	response.SuccessWithPagination(ctx, "获取日志成功", logs, page, pageSize, total)
}

// ClearLogs 清空日志
func (c *Controller) ClearLogs(ctx *gin.Context) {
	// 获取日志类型
	logType, err := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	if err != nil {
		response.BadRequest(ctx, "无效的日志类型")
		return
	}

	// 构建删除条件
	query := database.DB

	// 根据日志类型筛选
	if logType > 0 {
		query = query.Where("type = ?", logType)
	}

	// 应用ID筛选
	if appID := ctx.Query("app_id"); appID != "" {
		appIDInt, err := strconv.ParseUint(appID, 10, 32)
		if err == nil {
			query = query.Where("app_id = ?", uint(appIDInt))
		}
	}

	// 用户ID筛选
	if userID := ctx.Query("user_id"); userID != "" {
		userIDInt, err := strconv.ParseUint(userID, 10, 32)
		if err == nil {
			query = query.Where("user_id = ?", uint(userIDInt))
		}
	}

	// 清空日志
	result := query.Delete(&dbmodel.Logs{})
	if result.Error != nil {
		response.InternalServerError(ctx, "清空日志失败: "+result.Error.Error())
		return
	}

	// 返回清空成功响应
	response.Success(ctx, "清空日志成功", nil)
}
