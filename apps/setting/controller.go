package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/apps/setting/model"
	"github.com/skyle1995/DevE-Server/utils/response"
)

// Controller 处理系统设置相关的HTTP请求
type Controller struct {
	service *Service
}

// NewController 创建一个新的系统设置控制器实例
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// GetSettings 获取系统设置列表
func (c *Controller) GetSettings(ctx *gin.Context) {
	// 绑定请求参数
	var req model.GetSettingsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 调用服务层获取设置
	settings, err := c.service.GetSettings(req.Group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取设置失败: " + err.Error(),
		})
		return
	}

	// 返回设置列表
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    settings,
	})
}

// GetSettingByKey 根据键名获取系统设置
func (c *Controller) GetSettingByKey(ctx *gin.Context) {
	// 获取路径参数
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少键名参数",
		})
		return
	}

	// 调用服务层获取设置
	setting, err := c.service.GetSettingByKey(key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取设置失败: " + err.Error(),
		})
		return
	}

	// 返回设置
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    setting,
	})
}

// UpdateSetting 更新系统设置
func (c *Controller) UpdateSetting(ctx *gin.Context) {
	// 绑定请求参数
	var req model.UpdateSettingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层更新设置
	err := c.service.UpdateSetting(req.Key, req.Value)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新设置失败: " + err.Error(),
		})
		return
	}

	// 返回更新成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// CreateSetting 创建系统设置
func (c *Controller) CreateSetting(ctx *gin.Context) {
	// 绑定请求参数
	var req model.CreateSettingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 调用服务层创建设置
	err := c.service.CreateSetting(req.Key, req.Value, req.Group, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建设置失败: " + err.Error(),
		})
		return
	}

	// 返回创建成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

// DeleteSetting 删除系统设置
func (c *Controller) DeleteSetting(ctx *gin.Context) {
	// 获取路径参数
	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少键名参数",
		})
		return
	}

	// 调用服务层删除设置
	err := c.service.DeleteSetting(key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除设置失败: " + err.Error(),
		})
		return
	}

	// 返回删除成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// GetSiteInfo 获取站点信息
func (c *Controller) GetSiteInfo(ctx *gin.Context) {
	// 调用服务层获取站点信息
	siteInfo, err := c.service.GetSiteInfo()
	if err != nil {
		response.FailWithMessage("获取站点信息失败: "+err.Error(), ctx)
		return
	}

	// 构建响应数据
	data := model.SiteInfoResponse{
		Name:        siteInfo["name"],
		Subtitle:    siteInfo["subtitle"],
		Description: siteInfo["description"],
		Keywords:    siteInfo["keywords"],
		Copyright:   siteInfo["copyright"],
	}

	// 返回站点信息
	response.Success(ctx, "获取成功", data)
}
