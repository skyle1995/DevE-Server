package client

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/apps/client/model"
	"github.com/skyle1995/DevE-Server/utils/response"
)

// Controller 客户端控制器
type Controller struct {
	service *Service
}

// NewController 创建客户端控制器
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// ActivateCard 激活卡密
func (c *Controller) ActivateCard(ctx *gin.Context) {
	var req model.ActivateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文中获取应用信息
	app, exists := ctx.Get("app")
	if !exists {
		response.FailWithMessage("应用信息获取失败", ctx)
		return
	}

	// 调用服务层处理激活卡密
	res, err := c.service.ActivateCard(req, app)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithData(res, ctx)
}

// VerifyDevice 验证设备
func (c *Controller) VerifyDevice(ctx *gin.Context) {
	var req model.ActivateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文中获取应用信息
	app, exists := ctx.Get("app")
	if !exists {
		response.FailWithMessage("应用信息获取失败", ctx)
		return
	}

	// 调用服务层处理验证设备
	res, err := c.service.VerifyDevice(req, app)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithData(res, ctx)
}

// RebindCard 换绑卡密
func (c *Controller) RebindCard(ctx *gin.Context) {
	var req model.RebindCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文中获取应用信息
	app, exists := ctx.Get("app")
	if !exists {
		response.FailWithMessage("应用信息获取失败", ctx)
		return
	}

	// 调用服务层处理换绑卡密
	res, err := c.service.RebindCard(req, app)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithData(res, ctx)
}

// UnbindCard 解绑卡密
func (c *Controller) UnbindCard(ctx *gin.Context) {
	var req model.UnbindCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文中获取应用信息
	app, exists := ctx.Get("app")
	if !exists {
		response.FailWithMessage("应用信息获取失败", ctx)
		return
	}

	// 调用服务层处理解绑卡密
	res, err := c.service.UnbindCard(req, app)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithData(res, ctx)
}

// VerifyApp 验证应用
func (c *Controller) VerifyApp(ctx *gin.Context) {
	var req model.VerifyAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文中获取应用信息
	_, exists := ctx.Get("app")
	if !exists {
		response.FailWithMessage("应用信息获取失败", ctx)
		return
	}

	// 验证应用成功，直接返回成功响应
	res := model.VerifyAppResponse{
		Success: true,
		Message: "应用验证成功",
	}

	response.OkWithData(res, ctx)
}

// Heartbeat 处理客户端心跳
func (c *Controller) Heartbeat(ctx *gin.Context) {
	var req model.HeartbeatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 从上下文中获取应用信息
	app, exists := ctx.Get("app")
	if !exists {
		response.FailWithMessage("应用信息获取失败", ctx)
		return
	}

	// 调用服务层处理心跳
	res, err := c.service.Heartbeat(req, app)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithData(res, ctx)
}
