package apps

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/apps/app/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/middleware"
	"github.com/skyle1995/DevE-Server/utils/response"
)

// Controller 处理应用相关的HTTP请求
type Controller struct {
	service *Service
}

// NewController 创建一个新的应用控制器实例
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// CreateApp 创建新应用
func (c *Controller) CreateApp(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	// 绑定请求参数
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Version     string `json:"version"`
		DownloadUrl string `json:"download_url"`
		BillingMode int    `json:"billing_mode"`
		TrialAmount int    `json:"trial_amount"`
		AllowTrial  int    `json:"allow_trial"`
		PublicData  string `json:"public_data"`  // 公有数据
		PrivateData string `json:"private_data"` // 私有数据
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层创建应用
	app, err := c.service.CreateApp(
		userID.(uint),
		req.Name,
		req.Description,
		req.Version,
		req.DownloadUrl,
		req.BillingMode,
		req.TrialAmount,
		req.AllowTrial,
		req.PublicData,
		req.PrivateData,
	)

	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 转换为响应模型
	appResponse := model.NewAppResponse(app)

	// 返回创建成功响应
	response.Success(ctx, "应用创建成功", appResponse)
}

// GetAppList 获取应用列表
func (c *Controller) GetAppList(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// 调用服务层获取应用列表
	apps, total, err := c.service.GetAppList(userID.(uint), page, pageSize)
	if err != nil {
		response.InternalServerError(ctx, "获取应用列表失败: "+err.Error())
		return
	}

	// 转换为响应模型列表
	appResponses := model.NewAppResponseList(apps)

	// 返回应用列表
	response.SuccessWithPagination(ctx, "获取应用列表成功", appResponses, page, pageSize, total)
}

// GetAppByID 获取应用详情
func (c *Controller) GetAppByID(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	// 获取应用ID
	appID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的应用ID")
		return
	}

	// 调用服务层获取应用详情
	app, err := c.service.GetAppByID(userID.(uint), uint(appID))
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	// 转换为响应模型
	appResponse := model.NewAppResponseWithoutSecret(app)

	// 返回应用详情
	response.Success(ctx, "获取应用详情成功", appResponse)
}

// UpdateApp 更新应用信息
func (c *Controller) UpdateApp(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	// 获取应用ID
	appID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的应用ID")
		return
	}

	// 绑定请求参数
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
		DownloadUrl string `json:"download_url"`
		BillingMode int    `json:"billing_mode"`
		TrialAmount int    `json:"trial_amount"`
		AllowTrial  int    `json:"allow_trial"`
		Status      int    `json:"status"`
		PublicData  string `json:"public_data"`  // 公有数据
		PrivateData string `json:"private_data"` // 私有数据
	}

	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		response.BadRequest(ctx, "请求参数错误: "+bindErr.Error())
		return
	}

	// 调用服务层更新应用
	app, err := c.service.UpdateApp(
		userID.(uint),
		uint(appID),
		req.Name,
		req.Description,
		req.Version,
		req.DownloadUrl,
		req.BillingMode,
		req.TrialAmount,
		req.AllowTrial,
		req.Status,
		req.PublicData,
		req.PrivateData,
	)

	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 转换为响应模型
	appResponse := model.NewAppResponseWithoutSecret(app)

	// 返回更新成功响应
	response.Success(ctx, "应用更新成功", appResponse)
}

// DeleteApp 删除应用
func (c *Controller) DeleteApp(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	// 获取应用ID
	appID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的应用ID")
		return
	}

	// 调用服务层删除应用
	err = c.service.DeleteApp(userID.(uint), uint(appID))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 返回删除成功响应
	response.Success(ctx, "应用删除成功", nil)
}

// RegenerateAppSecret 重新生成应用密钥
func (c *Controller) RegenerateAppSecret(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	// 获取应用ID
	appID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的应用ID")
		return
	}

	// 调用服务层重新生成应用密钥
	app, err := c.service.RegenerateAppSecret(userID.(uint), uint(appID))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 转换为响应模型
	appResponse := model.NewAppResponse(app)

	// 返回重新生成密钥成功响应
	response.Success(ctx, "应用密钥重新生成成功", appResponse)
}

// SetupAppAPIRoutes 设置应用API路由
func SetupAppAPIRoutes(r *gin.Engine) {
	// 应用API路由组
	api := r.Group("/api/v1/client")
	api.Use(middleware.AppAuthMiddleware(), middleware.LogMiddleware(3)) // 使用LogTypeApp(3)类型的日志中间件
	{
		// 验证应用有效性
		api.GET("/verify", VerifyApp)

		// 需要验证试用期的路由组
		trial := api.Group("/")
		trial.Use(middleware.AppTrialMiddleware())
		{
			// 这里可以添加需要验证试用期的路由
			trial.GET("/status", GetAppStatus)
		}
	}
}

// VerifyApp 验证应用有效性
func VerifyApp(c *gin.Context) {
	// 从上下文中获取应用ID
	appID, exists := c.Get("app_id")
	if !exists {
		response.Unauthorized(c, "未授权的访问")
		return
	}

	// 查询应用
	var app dbmodel.App
	result := database.DB.First(&app, appID)
	if result.Error != nil {
		response.Unauthorized(c, "应用不存在")
		return
	}

	// 检查应用状态
	if app.Status != 1 {
		response.Forbidden(c, "应用已被禁用")
		return
	}

	// 计算试用期信息
	var trialInfo struct {
		HasTrial      bool      `json:"has_trial"`
		TrialEndTime  time.Time `json:"trial_end_time,omitempty"`
		TrialTimeLeft int64     `json:"trial_time_left,omitempty"` // 剩余小时数
	}

	if app.AllowTrial == 1 && app.TrialAmount > 0 {
		trialInfo.HasTrial = true
		// 根据计费模式处理试用额度
		var trialDuration time.Duration
		if app.BillingMode == 0 { // 时长计费模式
			trialDuration = time.Duration(app.TrialAmount) * time.Hour
		} else { // 点数计费模式
			// 点数计费模式下，假设每个点数对应1小时的使用时间
			trialDuration = time.Duration(app.TrialAmount) * time.Hour
		}
		trialInfo.TrialEndTime = app.CreatedAt.Add(trialDuration)
		trialInfo.TrialTimeLeft = int64(time.Until(trialInfo.TrialEndTime).Hours())
		if trialInfo.TrialTimeLeft < 0 {
			trialInfo.TrialTimeLeft = 0
		}
	}

	// 返回应用信息
	response.Success(c, "应用验证成功", gin.H{
		"app_id":       app.ID,
		"app_name":     app.Name,
		"version":      app.Version,
		"download_url": app.DownloadUrl,
		"trial_info":   trialInfo,
	})
}

// GetAppStatus 获取应用状态
func GetAppStatus(c *gin.Context) {
	// 从上下文中获取应用ID
	appID, exists := c.Get("app_id")
	if !exists {
		response.Unauthorized(c, "未授权的访问")
		return
	}

	// 查询应用
	var app dbmodel.App
	result := database.DB.First(&app, appID)
	if result.Error != nil {
		response.Unauthorized(c, "应用不存在")
		return
	}

	// 返回应用状态
	response.Success(c, "获取应用状态成功", gin.H{
		"app_id":       app.ID,
		"app_name":     app.Name,
		"status":       app.Status,
		"version":      app.Version,
		"download_url": app.DownloadUrl,
	})
}
