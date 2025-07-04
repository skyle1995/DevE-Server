package notice

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/apps/notice/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/response"
)

// Controller 通知控制器
type Controller struct {
	service *Service
}

// NewController 创建一个新的通知控制器实例
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// CreateNotice 创建通知
// @Summary 创建通知
// @Description 创建一个新的系统通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param data body model.CreateNoticeRequest true "通知信息"
// @Success 200 {object} response.Response{data=model.NoticeResponse}
// @Router /api/v1/notices [post]
func (c *Controller) CreateNotice(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.FailWithMessage("用户未登录", ctx)
		return
	}

	// 解析请求
	var req model.CreateNoticeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 创建通知
	notice, err := c.service.CreateNotice(req, userID.(uint))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 获取用户名
	var user dbmodel.User
	database.DB.Select("username").Where("id = ?", userID).First(&user)

	// 返回响应
	response.OkWithData(model.ConvertFromDBModel(*notice, user.Username), ctx)
}

// UpdateNotice 更新通知
// @Summary 更新通知
// @Description 更新指定ID的通知信息
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Param data body model.UpdateNoticeRequest true "通知信息"
// @Success 200 {object} response.Response{data=model.NoticeResponse}
// @Router /api/v1/notices/{id} [put]
func (c *Controller) UpdateNotice(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.FailWithMessage("用户未登录", ctx)
		return
	}

	// 获取通知ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的通知ID", ctx)
		return
	}

	// 解析请求
	var req model.UpdateNoticeRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		response.FailWithMessage("参数错误: "+bindErr.Error(), ctx)
		return
	}

	// 更新通知
	notice, err := c.service.UpdateNotice(uint(id), req, userID.(uint))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 获取用户名
	var user dbmodel.User
	database.DB.Select("username").Where("id = ?", notice.UserID).First(&user)

	// 返回响应
	response.OkWithData(model.ConvertFromDBModel(*notice, user.Username), ctx)
}

// DeleteNotice 删除通知
// @Summary 删除通知
// @Description 删除指定ID的通知
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response
// @Router /api/v1/notices/{id} [delete]
func (c *Controller) DeleteNotice(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.FailWithMessage("用户未登录", ctx)
		return
	}

	// 获取通知ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的通知ID", ctx)
		return
	}

	// 删除通知
	err = c.service.DeleteNotice(uint(id), userID.(uint))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 返回响应
	response.OkWithMessage("删除成功", ctx)
}

// GetNotice 获取通知详情
// @Summary 获取通知详情
// @Description 获取指定ID的通知详情
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} response.Response{data=model.NoticeResponse}
// @Router /api/v1/notices/{id} [get]
func (c *Controller) GetNotice(ctx *gin.Context) {
	// 获取通知ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的通知ID", ctx)
		return
	}

	// 获取通知
	notice, err := c.service.GetNoticeByID(uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 获取用户名
	var user dbmodel.User
	database.DB.Select("username").Where("id = ?", notice.UserID).First(&user)

	// 返回响应
	response.OkWithData(model.ConvertFromDBModel(*notice, user.Username), ctx)
}

// GetNoticeList 获取通知列表
// @Summary 获取通知列表
// @Description 获取通知列表，支持分页和筛选
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param title query string false "标题（模糊查询）"
// @Param level query int false "等级"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=model.NoticeListResponse}
// @Router /api/v1/notices [get]
func (c *Controller) GetNoticeList(ctx *gin.Context) {
	// 解析请求
	var req model.GetNoticeListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 获取通知列表
	notices, total, err := c.service.GetNoticeList(req, false)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 构建响应
	list := make([]model.NoticeResponse, 0, len(notices))
	for _, notice := range notices {
		// 获取用户名
		var user dbmodel.User
		database.DB.Select("username").Where("id = ?", notice.UserID).First(&user)

		list = append(list, model.ConvertFromDBModel(notice, user.Username))
	}

	// 返回响应
	response.OkWithData(model.NoticeListResponse{
		Total: int(total),
		List:  list,
	}, ctx)
}

// GetActiveNoticeList 获取活动通知列表
// @Summary 获取活动通知列表
// @Description 获取当前有效的通知列表，支持分页
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} response.Response{data=model.NoticeListResponse}
// @Router /api/v1/notices/active [get]
func (c *Controller) GetActiveNoticeList(ctx *gin.Context) {
	// 解析请求
	var req model.GetNoticeListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	// 获取活动通知列表
	notices, total, err := c.service.GetNoticeList(req, true)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 构建响应
	list := make([]model.NoticeResponse, 0, len(notices))
	for _, notice := range notices {
		// 获取用户名
		var user dbmodel.User
		database.DB.Select("username").Where("id = ?", notice.UserID).First(&user)

		list = append(list, model.ConvertFromDBModel(notice, user.Username))
	}

	// 返回响应
	response.OkWithData(model.NoticeListResponse{
		Total: int(total),
		List:  list,
	}, ctx)
}

// GetDashboardNotices 获取仪表盘公告列表
// @Summary 获取仪表盘公告列表
// @Description 获取仪表盘显示的公告列表
// @Tags 通知管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]model.NoticeResponse}
// @Router /api/v1/notices/dashboard [get]
func (c *Controller) GetDashboardNotices(ctx *gin.Context) {
	// 构建请求参数
	req := model.GetNoticeListRequest{
		Page:     1,
		PageSize: 10, // 限制返回10条公告
	}

	// 获取活动通知列表
	notices, _, err := c.service.GetNoticeList(req, true)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 构建响应
	list := make([]model.NoticeResponse, 0, len(notices))
	for _, notice := range notices {
		// 获取用户名
		var user dbmodel.User
		database.DB.Select("username").Where("id = ?", notice.UserID).First(&user)

		list = append(list, model.ConvertFromDBModel(notice, user.Username))
	}

	// 返回响应
	response.OkWithData(list, ctx)
}

// UpdateNoticeStatus 更新通知状态
// @Summary 更新通知状态
// @Description 更新指定ID的通知状态
// @Tags 通知管理
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Param data body model.UpdateNoticeStatusRequest true "状态信息"
// @Success 200 {object} response.Response{data=model.NoticeResponse}
// @Router /api/v1/notices/{id}/status [put]
func (c *Controller) UpdateNoticeStatus(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.FailWithMessage("用户未登录", ctx)
		return
	}

	// 获取通知ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的通知ID", ctx)
		return
	}

	// 解析请求
	var req model.UpdateNoticeStatusRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		response.FailWithMessage("参数错误: "+bindErr.Error(), ctx)
		return
	}

	// 更新通知状态
	notice, err := c.service.UpdateNoticeStatus(uint(id), req.Status, userID.(uint))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 获取用户名
	var user dbmodel.User
	database.DB.Select("username").Where("id = ?", notice.UserID).First(&user)

	// 返回响应
	response.OkWithData(model.ConvertFromDBModel(*notice, user.Username), ctx)
}
