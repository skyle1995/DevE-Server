package card

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skyle1995/DevE-Server/apps/card/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/response"
	"gorm.io/gorm"
)

// Controller 卡密控制器
type Controller struct {
	service *Service
}

// NewController 创建一个新的卡密控制器实例
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// GetCardTypes 获取卡密类型列表
// @Summary 获取卡密类型列表
// @Description 获取当前用户创建的卡密类型列表
// @Tags 用户API
// @Accept json
// @Produce json
// @Param app_id query int false "应用ID"
// @Success 200 {object} response.Response{data=model.CardTypeListResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/types [get]
func (c *Controller) GetCardTypes(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 解析应用ID参数
	appIDStr := ctx.Query("app_id")
	var appID *int
	if appIDStr != "" {
		id, err := strconv.Atoi(appIDStr)
		if err != nil {
			response.FailWithMessage("应用ID格式错误", ctx)
			return
		}
		appID = &id
	}

	// 查询卡密类型列表
	var cardTypes []dbmodel.CardType
	db := database.DB.Where("user_id = ?", userID)
	if appID != nil {
		db = db.Where("app_id = ?", *appID)
	}
	result := db.Find(&cardTypes)
	if result.Error != nil {
		response.FailWithMessage("查询卡密类型失败: "+result.Error.Error(), ctx)
		return
	}

	// 构建响应
	response.OkWithData(model.FromCardTypes(cardTypes), ctx)
}

// CreateCardType 创建卡密类型
// @Summary 创建卡密类型
// @Description 创建新的卡密类型
// @Tags 用户API
// @Accept json
// @Produce json
// @Param request body model.CreateCardTypeRequest true "创建卡密类型请求"
// @Success 200 {object} response.Response{data=model.CardTypeResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/types [post]
func (c *Controller) CreateCardType(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 绑定请求参数
	var req model.CreateCardTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), ctx)
		return
	}

	// 创建卡密类型
	cardType := dbmodel.CardType{
		Name:                  req.Name,
		Duration:              req.Duration,
		TimeUnit:              req.TimeUnit,
		AppID:                 uint(req.AppID),
		Status:                req.Status,
		DefaultMaxRebindCount: req.DefaultMaxRebindCount,
		DefaultMaxUnbindCount: req.DefaultMaxUnbindCount,
		UserID:                userID.(int),
	}

	result := database.DB.Create(&cardType)
	if result.Error != nil {
		response.FailWithMessage("创建卡密类型失败: "+result.Error.Error(), ctx)
		return
	}

	// 构建响应
	response.OkWithData(model.FromCardType(cardType), ctx)
}

// UpdateCardType 更新卡密类型
// @Summary 更新卡密类型
// @Description 更新卡密类型信息
// @Tags 用户API
// @Accept json
// @Produce json
// @Param id path int true "卡密类型ID"
// @Param request body model.UpdateCardTypeRequest true "更新卡密类型请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/types/{id} [put]
func (c *Controller) UpdateCardType(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 解析卡密类型ID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.FailWithMessage("卡密类型ID格式错误", ctx)
		return
	}

	// 绑定请求参数
	var req model.UpdateCardTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), ctx)
		return
	}

	// 查询卡密类型
	var cardType dbmodel.CardType
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&cardType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.FailWithMessage("卡密类型不存在或无权限修改", ctx)
			return
		}
		response.FailWithMessage("查询卡密类型失败: "+result.Error.Error(), ctx)
		return
	}

	// 更新卡密类型
	if req.Name != "" {
		cardType.Name = req.Name
	}
	if req.Duration > 0 {
		cardType.Duration = req.Duration
	}
	if req.TimeUnit != "" {
		cardType.TimeUnit = req.TimeUnit
	}
	if req.Status >= 0 {
		cardType.Status = req.Status
	}

	result = database.DB.Save(&cardType)
	if result.Error != nil {
		response.FailWithMessage("更新卡密类型失败: "+result.Error.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// DeleteCardType 删除卡密类型
// @Summary 删除卡密类型
// @Description 删除卡密类型
// @Tags 用户API
// @Accept json
// @Produce json
// @Param id path int true "卡密类型ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/types/{id} [delete]
func (c *Controller) DeleteCardType(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 解析卡密类型ID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.FailWithMessage("卡密类型ID格式错误", ctx)
		return
	}

	// 查询卡密类型
	var cardType dbmodel.CardType
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&cardType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.FailWithMessage("卡密类型不存在或无权限删除", ctx)
			return
		}
		response.FailWithMessage("查询卡密类型失败: "+result.Error.Error(), ctx)
		return
	}

	// 删除卡密类型
	result = database.DB.Delete(&cardType)
	if result.Error != nil {
		response.FailWithMessage("删除卡密类型失败: "+result.Error.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// GetCards 获取卡密列表
// @Summary 获取卡密列表
// @Description 获取当前用户创建的卡密列表
// @Tags 用户API
// @Accept json
// @Produce json
// @Param app_id query int false "应用ID"
// @Param type_id query int false "卡密类型ID"
// @Param status query int false "卡密状态"
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} response.Response{data=model.CardListResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/cards [get]
func (c *Controller) GetCards(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 解析查询参数
	var req model.GetCardListRequest
	// 应用ID
	appIDStr := ctx.Query("app_id")
	if appIDStr != "" {
		id, err := strconv.Atoi(appIDStr)
		if err != nil {
			response.FailWithMessage("应用ID格式错误", ctx)
			return
		}
		req.AppID = id
	}
	// 卡密类型ID
	typeIDStr := ctx.Query("type_id")
	if typeIDStr != "" {
		id, err := strconv.Atoi(typeIDStr)
		if err != nil {
			response.FailWithMessage("卡密类型ID格式错误", ctx)
			return
		}
		req.TypeID = id
	}
	// 卡密状态
	statusStr := ctx.Query("status")
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			response.FailWithMessage("卡密状态格式错误", ctx)
			return
		}
		req.Status = status
	}
	// 分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "20")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	req.Page = page
	req.PageSize = pageSize

	// 查询卡密列表
	var cards []dbmodel.Card
	db := database.DB.Preload("CardType").Where("user_id = ?", userID)
	if req.AppID > 0 {
		db = db.Where("app_id = ?", req.AppID)
	}
	if req.TypeID > 0 {
		db = db.Where("type_id = ?", req.TypeID)
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 计算总数
	var total int64
	db.Count(&total)

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	result := db.Offset(offset).Limit(req.PageSize).Find(&cards)
	if result.Error != nil {
		response.FailWithMessage("查询卡密列表失败: "+result.Error.Error(), ctx)
		return
	}

	// 构建响应
	response.OkWithData(model.CardListResponse{
		Total: int(total),
		Items: model.FromCards(cards, true),
	}, ctx)
}

// GenerateCards 生成卡密
// @Summary 生成卡密
// @Description 批量生成卡密
// @Tags 用户API
// @Accept json
// @Produce json
// @Param request body model.GenerateCardRequest true "生成卡密请求"
// @Success 200 {object} response.Response{data=model.GenerateCardResponse} "生成成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/generate [post]
func (c *Controller) GenerateCards(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 绑定请求参数
	var req model.GenerateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), ctx)
		return
	}

	// 验证卡密长度是否符合最低要求
	if req.KeyLength != nil && *req.KeyLength < 8 {
		response.FailWithMessage("卡密长度不能小于8", ctx)
		return
	}

	// 验证卡密类型是否存在且属于当前用户
	var cardType dbmodel.CardType
	result := database.DB.Where("id = ? AND user_id = ?", req.TypeID, userID).First(&cardType)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.FailWithMessage("卡密类型不存在或无权限使用", ctx)
			return
		}
		response.FailWithMessage("查询卡密类型失败: "+result.Error.Error(), ctx)
		return
	}

	// 生成卡密
	cards, err := c.service.GenerateCards(req, userID.(int))
	if err != nil {
		response.FailWithMessage("生成卡密失败: "+err.Error(), ctx)
		return
	}

	// 构建响应
	response.OkWithData(model.GenerateCardResponse{
		Success: true,
		Count:   len(cards),
		Cards:   model.FromCards(cards, true),
	}, ctx)
}

// UpdateCard 更新卡密
// @Summary 更新卡密
// @Description 更新卡密信息
// @Tags 用户API
// @Accept json
// @Produce json
// @Param id path int true "卡密ID"
// @Param request body model.UpdateCardRequest true "更新卡密请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/cards/{id} [put]
func (c *Controller) UpdateCard(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 解析卡密ID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.FailWithMessage("卡密ID格式错误", ctx)
		return
	}

	// 绑定请求参数
	var req model.UpdateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+err.Error(), ctx)
		return
	}

	// 查询卡密
	var card dbmodel.Card
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.FailWithMessage("卡密不存在或无权限修改", ctx)
			return
		}
		response.FailWithMessage("查询卡密失败: "+result.Error.Error(), ctx)
		return
	}

	// 更新卡密
	if req.TypeID > 0 {
		// 验证卡密类型是否存在且属于当前用户
		var cardType dbmodel.CardType
		result = database.DB.Where("id = ? AND user_id = ?", req.TypeID, userID).First(&cardType)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				response.FailWithMessage("卡密类型不存在或无权限使用", ctx)
				return
			}
			response.FailWithMessage("查询卡密类型失败: "+result.Error.Error(), ctx)
			return
		}
		card.TypeID = uint(req.TypeID)
	}

	if req.Status >= 0 {
		card.Status = req.Status
	}

	result = database.DB.Save(&card)
	if result.Error != nil {
		response.FailWithMessage("更新卡密失败: "+result.Error.Error(), ctx)
		return
	}

	response.Ok(ctx)
}

// DeleteCard 删除卡密
// @Summary 删除卡密
// @Description 删除卡密
// @Tags 用户API
// @Accept json
// @Produce json
// @Param id path int true "卡密ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/card/cards/{id} [delete]
func (c *Controller) DeleteCard(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.FailWithMessage("未找到用户信息", ctx)
		return
	}

	// 解析卡密ID
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.FailWithMessage("卡密ID格式错误", ctx)
		return
	}

	// 查询卡密
	var card dbmodel.Card
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.FailWithMessage("卡密不存在或无权限删除", ctx)
			return
		}
		response.FailWithMessage("查询卡密失败: "+result.Error.Error(), ctx)
		return
	}

	// 删除卡密
	result = database.DB.Delete(&card)
	if result.Error != nil {
		response.FailWithMessage("删除卡密失败: "+result.Error.Error(), ctx)
		return
	}

	response.Ok(ctx)
}
