package notice

import (
	"errors"
	"time"

	"github.com/skyle1995/DevE-Server/apps/notice/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"gorm.io/gorm"
)

// Service 提供通知相关的服务
type Service struct{}

// NewService 创建一个新的通知服务实例
func NewService() *Service {
	return &Service{}
}

// CreateNotice 创建通知
// @param req 创建通知请求
// @param userID 当前用户ID
// @return 创建的通知和错误信息
func (s *Service) CreateNotice(req model.CreateNoticeRequest, userID uint) (*dbmodel.Notice, error) {
	// 创建通知
	notice := dbmodel.Notice{
		Title:     req.Title,
		Content:   req.Content,
		Level:     req.Level,
		Status:    req.Status,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		UserID:    userID,
	}

	// 保存到数据库
	result := database.DB.Create(&notice)
	if result.Error != nil {
		return nil, errors.New("创建通知失败: " + result.Error.Error())
	}

	return &notice, nil
}

// UpdateNotice 更新通知
// @param id 通知ID
// @param req 更新通知请求
// @param userID 当前用户ID
// @return 更新后的通知和错误信息
func (s *Service) UpdateNotice(id uint, req model.UpdateNoticeRequest, userID uint) (*dbmodel.Notice, error) {
	// 查询通知
	var notice dbmodel.Notice
	result := database.DB.Where("id = ?", id).First(&notice)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("通知不存在")
		}
		return nil, errors.New("查询通知失败: " + result.Error.Error())
	}

	// 检查权限（只有管理员或发布者可以更新）
	if notice.UserID != userID {
		// 检查当前用户是否为管理员
		var user dbmodel.User
		userResult := database.DB.Where("id = ?", userID).First(&user)
		if userResult.Error != nil || user.Role != 1 {
			return nil, errors.New("无权限更新此通知")
		}
	}

	// 更新通知
	if req.Title != "" {
		notice.Title = req.Title
	}
	if req.Content != "" {
		notice.Content = req.Content
	}
	if req.Level != nil {
		notice.Level = *req.Level
	}
	if req.Status != nil {
		notice.Status = *req.Status
	}
	if req.StartTime != nil {
		notice.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		notice.EndTime = req.EndTime
	}

	// 保存到数据库
	result = database.DB.Save(&notice)
	if result.Error != nil {
		return nil, errors.New("更新通知失败: " + result.Error.Error())
	}

	return &notice, nil
}

// DeleteNotice 删除通知
// @param id 通知ID
// @param userID 当前用户ID
// @return 错误信息
func (s *Service) DeleteNotice(id uint, userID uint) error {
	// 查询通知
	var notice dbmodel.Notice
	result := database.DB.Where("id = ?", id).First(&notice)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("通知不存在")
		}
		return errors.New("查询通知失败: " + result.Error.Error())
	}

	// 检查权限（只有管理员或发布者可以删除）
	if notice.UserID != userID {
		// 检查当前用户是否为管理员
		var user dbmodel.User
		userResult := database.DB.Where("id = ?", userID).First(&user)
		if userResult.Error != nil || user.Role != 1 {
			return errors.New("无权限删除此通知")
		}
	}

	// 删除通知
	result = database.DB.Delete(&notice)
	if result.Error != nil {
		return errors.New("删除通知失败: " + result.Error.Error())
	}

	return nil
}

// GetNoticeByID 根据ID获取通知
// @param id 通知ID
// @return 通知和错误信息
func (s *Service) GetNoticeByID(id uint) (*dbmodel.Notice, error) {
	// 查询通知
	var notice dbmodel.Notice
	result := database.DB.Where("id = ?", id).First(&notice)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("通知不存在")
		}
		return nil, errors.New("查询通知失败: " + result.Error.Error())
	}

	return &notice, nil
}

// GetNoticeList 获取通知列表
// @param req 获取通知列表请求
// @param onlyActive 是否只获取启用状态的通知
// @return 通知列表、总数和错误信息
func (s *Service) GetNoticeList(req model.GetNoticeListRequest, onlyActive bool) ([]dbmodel.Notice, int64, error) {
	// 构建查询
	query := database.DB.Model(&dbmodel.Notice{})

	// 添加查询条件
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.Level != nil {
		query = query.Where("level = ?", *req.Level)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	} else if onlyActive {
		query = query.Where("status = ?", 1) // 只获取启用状态的通知
	}

	// 如果只获取有效的通知，添加时间条件
	if onlyActive {
		now := time.Now()
		query = query.Where(
			"(start_time IS NULL OR start_time <= ?) AND (end_time IS NULL OR end_time >= ?)",
			now, now,
		)
	}

	// 获取总数
	var total int64
	result := query.Count(&total)
	if result.Error != nil {
		return nil, 0, errors.New("获取通知总数失败: " + result.Error.Error())
	}

	// 分页查询
	var notices []dbmodel.Notice
	result = query.Order("created_at DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&notices)
	if result.Error != nil {
		return nil, 0, errors.New("获取通知列表失败: " + result.Error.Error())
	}

	return notices, total, nil
}

// UpdateNoticeStatus 更新通知状态
// @param id 通知ID
// @param status 状态
// @param userID 当前用户ID
// @return 更新后的通知和错误信息
func (s *Service) UpdateNoticeStatus(id uint, status int, userID uint) (*dbmodel.Notice, error) {
	// 查询通知
	var notice dbmodel.Notice
	result := database.DB.Where("id = ?", id).First(&notice)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("通知不存在")
		}
		return nil, errors.New("查询通知失败: " + result.Error.Error())
	}

	// 检查权限（只有管理员或发布者可以更新状态）
	if notice.UserID != userID {
		// 检查当前用户是否为管理员
		var user dbmodel.User
		userResult := database.DB.Where("id = ?", userID).First(&user)
		if userResult.Error != nil || user.Role != 1 {
			return nil, errors.New("无权限更新此通知状态")
		}
	}

	// 更新状态
	notice.Status = status

	// 保存到数据库
	result = database.DB.Save(&notice)
	if result.Error != nil {
		return nil, errors.New("更新通知状态失败: " + result.Error.Error())
	}

	return &notice, nil
}
