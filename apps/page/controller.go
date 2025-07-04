package page

import (
	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/utils/response"
)

// Controller 页面控制器结构体
type Controller struct {
	service *Service
}

// NewController 创建页面控制器实例
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// GetMenu 获取菜单数据
// @Summary 获取菜单数据
// @Description 根据用户角色获取对应的菜单数据
// @Tags 页面管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=model.MenuResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/menu [get]
func (c *Controller) GetMenu(ctx *gin.Context) {
	// 从上下文中获取用户角色
	roleValue, exists := ctx.Get("role")
	if !exists {
		response.FailWithMessage("获取用户角色失败", ctx)
		return
	}

	role, ok := roleValue.(int)
	if !ok {
		response.FailWithMessage("用户角色类型错误", ctx)
		return
	}

	// 获取菜单数据
	menuResp, err := c.service.GetMenu(role)
	if err != nil {
		response.FailWithMessage("获取菜单数据失败", ctx)
		return
	}

	response.OkWithData(menuResp, ctx)
}
