package client

import (
	"errors"
	"time"

	"github.com/skyle1995/DevE-Server/apps/client/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/timeutil"
	"gorm.io/gorm"
)

// Service 客户端服务
type Service struct {
	db *gorm.DB
}

// NewService 创建客户端服务
func NewService() *Service {
	return &Service{
		db: database.DB,
	}
}

// ActivateCard 激活卡密
func (s *Service) ActivateCard(req model.ActivateCardRequest, app interface{}) (*model.ActivateCardResponse, error) {
	// 类型断言获取应用信息
	appInfo, ok := app.(dbmodel.App)
	if !ok {
		return nil, errors.New("应用信息类型错误")
	}

	// 查询卡密
	var card dbmodel.Card
	result := s.db.Where("card_no = ? AND app_id = ?", req.CardNo, appInfo.ID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("卡密不存在")
		}
		return nil, result.Error
	}

	// 检查卡密状态
	if card.Status == 2 { // 已禁用
		return nil, errors.New("卡密已被禁用")
	}

	// 检查卡密是否已过期
	if card.ExpireAt != nil && card.ExpireAt.Before(time.Now()) {
		return nil, errors.New("卡密已过期")
	}

	// 查询卡密类型
	var cardType dbmodel.CardType
	result = s.db.First(&cardType, card.TypeID)
	if result.Error != nil {
		return nil, errors.New("卡密类型不存在")
	}

	// 查询应用设置
	var appSetting dbmodel.App
	result = s.db.Where("app_id = ?", appInfo.ID).First(&appSetting)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("获取应用设置失败")
	}

	// 如果卡密已激活，检查设备ID是否匹配
	if card.Status == 1 { // 已激活
		if card.DeviceID != nil && *card.DeviceID != req.DeviceID {
			return nil, errors.New("卡密已绑定其他设备")
		}

		// 返回卡密信息
		return &model.ActivateCardResponse{
			CardNo:        card.CardNo,
			Status:        card.Status,
			Activated:     true,
			ExpireAt:      card.ExpireAt,
			DeviceBinding: card.DeviceID != nil,
			BindingInfo:   "已绑定当前设备",
			BindCount:     card.RebindCount,
			MaxBindCount:  cardType.MaxBindCount,
			CanRebind:     (appSetting.BindPermission == 1 || appSetting.BindPermission == 3) && (cardType.MaxBindCount == 0 || card.RebindCount < cardType.MaxBindCount),
			CanUnbind:     appSetting.BindPermission == 2 || appSetting.BindPermission == 3,
			Message:       "卡密已激活",
		}, nil
	}

	// 查询或创建设备
	var device dbmodel.Device
	result = s.db.Where("device_id = ? AND app_id = ?", req.DeviceID, appInfo.ID).First(&device)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 创建新设备
		device = dbmodel.Device{
			DeviceID:   req.DeviceID,
			AppID:      appInfo.ID,
			DeviceInfo: req.DeviceInfo,
			Status:     1, // 正常状态
			LastActive: time.Now(),
		}
		result = s.db.Create(&device)
		if result.Error != nil {
			return nil, errors.New("创建设备记录失败")
		}
	} else if result.Error != nil {
		return nil, errors.New("查询设备信息失败")
	} else {
		// 更新设备信息
		device.DeviceInfo = req.DeviceInfo
		device.LastActive = time.Now()
		s.db.Save(&device)
	}

	// 检查设备状态
	if device.Status != 1 {
		return nil, errors.New("设备已被禁用")
	}

	// 计算过期时间
	var expireAt *time.Time
	if cardType.ValidDays > 0 {
		now := time.Now()
		expire := now.AddDate(0, 0, cardType.ValidDays)
		expireAt = &expire
	}

	// 更新卡密信息
	card.Status = 1 // 已激活
	card.DeviceID = &req.DeviceID
	now := time.Now()
	card.ActivateAt = &now
	card.ExpireAt = expireAt
	card.RebindCount = 0 // 初始化换绑次数

	result = s.db.Save(&card)
	if result.Error != nil {
		return nil, errors.New("更新卡密信息失败")
	}

	// 返回激活成功响应
	return &model.ActivateCardResponse{
		CardNo:        card.CardNo,
		Status:        card.Status,
		Activated:     true,
		ExpireAt:      card.ExpireAt,
		DeviceBinding: true,
		BindingInfo:   "已绑定当前设备",
		BindCount:     card.RebindCount,
		MaxBindCount:  cardType.MaxBindCount,
		CanRebind:     (appSetting.BindPermission == 1 || appSetting.BindPermission == 3) && (cardType.MaxBindCount == 0 || card.RebindCount < cardType.MaxBindCount),
		CanUnbind:     appSetting.BindPermission == 2 || appSetting.BindPermission == 3,
		Message:       "卡密激活成功",
	}, nil
}

// VerifyDevice 验证设备
func (s *Service) VerifyDevice(req model.ActivateCardRequest, app interface{}) (*model.VerifyDeviceResponse, error) {
	// 类型断言获取应用信息
	appInfo, ok := app.(dbmodel.App)
	if !ok {
		return nil, errors.New("应用信息类型错误")
	}

	// 查询设备
	var device dbmodel.Device
	result := s.db.Where("device_id = ? AND app_id = ?", req.DeviceID, appInfo.ID).First(&device)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("设备未注册")
		}
		return nil, errors.New("查询设备信息失败")
	}

	// 检查设备状态
	if device.Status != 1 {
		return nil, errors.New("设备已被禁用")
	}

	// 更新设备活跃时间
	device.LastActive = time.Now()
	s.db.Save(&device)

	// 查询关联的卡密
	var card dbmodel.Card
	result = s.db.Where("device_id = ? AND app_id = ? AND status = 1", req.DeviceID, appInfo.ID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到关联的卡密")
		}
		return nil, errors.New("查询卡密信息失败")
	}

	// 检查卡密是否过期
	if card.ExpireAt != nil && card.ExpireAt.Before(time.Now()) {
		return nil, errors.New("卡密已过期")
	}

	// 更新卡密在线状态和心跳时间
	now := time.Now()
	card.IsOnline = 1 // 设置为在线
	card.LastHeartbeat = &now
	s.db.Save(&card)

	// 计算剩余天数
	remainDays := 0
	if card.ExpireAt != nil {
		remainDays = timeutil.DaysBetween(time.Now(), *card.ExpireAt)
	}

	// 返回验证成功响应
	return &model.VerifyDeviceResponse{
		Success:    true,
		CardNo:     card.CardNo,
		ExpireAt:   card.ExpireAt,
		RemainDays: remainDays,
		Message:    "设备验证成功",
	}, nil
}

// RebindCard 换绑卡密
func (s *Service) RebindCard(req model.RebindCardRequest, app interface{}) (*model.RebindCardResponse, error) {
	// 类型断言获取应用信息
	appInfo, ok := app.(dbmodel.App)
	if !ok {
		return nil, errors.New("应用信息类型错误")
	}

	// 查询卡密
	var card dbmodel.Card
	result := s.db.Where("card_no = ? AND app_id = ?", req.CardNo, appInfo.ID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("卡密不存在")
		}
		return nil, result.Error
	}

	// 检查卡密状态
	if card.Status != 1 {
		return nil, errors.New("卡密未激活")
	}

	// 检查卡密是否已过期
	if card.ExpireAt != nil && card.ExpireAt.Before(time.Now()) {
		return nil, errors.New("卡密已过期")
	}

	// 检查卡密是否已禁用
	if card.Status == 2 {
		return nil, errors.New("卡密已被禁用")
	}

	// 检查卡密是否已绑定设备
	if card.DeviceID == nil {
		return nil, errors.New("卡密未绑定设备")
	}

	// 查询应用设置
	var appSetting dbmodel.App
	result = s.db.Where("app_id = ?", appInfo.ID).First(&appSetting)
	if result.Error != nil {
		return nil, errors.New("获取应用设置失败")
	}

	// 检查应用是否允许换绑
	if appSetting.BindPermission != 1 && appSetting.BindPermission != 3 {
		return nil, errors.New("应用不允许换绑")
	}

	// 查询卡密类型
	var cardType dbmodel.CardType
	result = s.db.First(&cardType, card.TypeID)
	if result.Error != nil {
		return nil, errors.New("卡密类型不存在")
	}

	// 检查换绑次数限制
	if cardType.MaxBindCount > 0 && card.RebindCount >= cardType.MaxBindCount {
		return nil, errors.New("已达到最大换绑次数限制")
	}

	// 查询或创建新设备
	var device dbmodel.Device
	result = s.db.Where("device_id = ? AND app_id = ?", req.DeviceID, appInfo.ID).First(&device)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 创建新设备
		device = dbmodel.Device{
			DeviceID:   req.DeviceID,
			AppID:      appInfo.ID,
			DeviceInfo: req.DeviceInfo,
			Status:     1, // 正常状态
			LastActive: time.Now(),
		}
		result = s.db.Create(&device)
		if result.Error != nil {
			return nil, errors.New("创建设备记录失败")
		}
	} else if result.Error != nil {
		return nil, errors.New("查询设备信息失败")
	} else {
		// 更新设备信息
		device.DeviceInfo = req.DeviceInfo
		device.LastActive = time.Now()
		s.db.Save(&device)
	}

	// 检查设备状态
	if device.Status != 1 {
		return nil, errors.New("设备已被禁用")
	}

	// 更新卡密信息
	card.DeviceID = &req.DeviceID
	card.RebindCount++

	result = s.db.Save(&card)
	if result.Error != nil {
		return nil, errors.New("更新卡密信息失败")
	}

	// 返回换绑成功响应
	return &model.RebindCardResponse{
		Success:  true,
		CardNo:   card.CardNo,
		DeviceID: req.DeviceID,
		Message:  "卡密换绑成功",
	}, nil
}

// Heartbeat 处理客户端心跳
func (s *Service) Heartbeat(req model.HeartbeatRequest, app interface{}) (*model.HeartbeatResponse, error) {
	// 类型断言获取应用信息
	appInfo, ok := app.(dbmodel.App)
	if !ok {
		return nil, errors.New("应用信息类型错误")
	}

	// 查询卡密
	var card dbmodel.Card
	result := s.db.Where("card_no = ? AND app_id = ? AND status = 1", req.CardNo, appInfo.ID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("卡密不存在或未激活")
		}
		return nil, errors.New("查询卡密信息失败")
	}

	// 检查设备ID是否匹配
	if card.DeviceID == nil || *card.DeviceID != req.DeviceID {
		return nil, errors.New("设备ID不匹配")
	}

	// 检查卡密是否过期
	if card.ExpireAt != nil && card.ExpireAt.Before(time.Now()) {
		// 更新卡密状态为已过期
		card.Status = 2   // 已过期
		card.IsOnline = 0 // 设置为离线
		s.db.Save(&card)
		return nil, errors.New("卡密已过期")
	}

	// 更新卡密在线状态和心跳时间
	now := time.Now()
	card.IsOnline = 1 // 设置为在线
	card.LastHeartbeat = &now
	s.db.Save(&card)

	// 更新设备活跃时间
	var device dbmodel.Device
	result = s.db.Where("device_id = ? AND app_id = ?", req.DeviceID, appInfo.ID).First(&device)
	if result.Error == nil {
		device.LastActive = now
		s.db.Save(&device)
	}

	// 计算剩余天数
	remainDays := 0
	if card.ExpireAt != nil {
		remainDays = int(card.ExpireAt.Sub(now).Hours() / 24)
		if remainDays < 0 {
			remainDays = 0
		}
	}

	// 返回心跳响应
	return &model.HeartbeatResponse{
		Success:    true,
		CardNo:     card.CardNo,
		IsOnline:   true,
		ExpireAt:   card.ExpireAt,
		RemainDays: remainDays,
		Message:    "心跳成功",
	}, nil
}

// UnbindCard 解绑卡密
func (s *Service) UnbindCard(req model.UnbindCardRequest, app interface{}) (*model.UnbindCardResponse, error) {
	// 类型断言获取应用信息
	appInfo, ok := app.(dbmodel.App)
	if !ok {
		return nil, errors.New("应用信息类型错误")
	}

	// 查询卡密
	var card dbmodel.Card
	result := s.db.Where("card_no = ? AND app_id = ?", req.CardNo, appInfo.ID).First(&card)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("卡密不存在")
		}
		return nil, result.Error
	}

	// 检查卡密状态
	if card.Status != 1 {
		return nil, errors.New("卡密未激活")
	}

	// 检查卡密是否已过期
	if card.ExpireAt != nil && card.ExpireAt.Before(time.Now()) {
		return nil, errors.New("卡密已过期")
	}

	// 检查卡密是否已禁用
	if card.Status == 2 {
		return nil, errors.New("卡密已被禁用")
	}

	// 检查卡密是否已绑定设备
	if card.DeviceID == nil {
		return nil, errors.New("卡密未绑定设备")
	}

	// 查询应用设置
	var appSetting dbmodel.App
	result = s.db.Where("app_id = ?", appInfo.ID).First(&appSetting)
	if result.Error != nil {
		return nil, errors.New("获取应用设置失败")
	}

	// 检查应用是否允许解绑
	if appSetting.BindPermission != 2 && appSetting.BindPermission != 3 {
		return nil, errors.New("应用不允许解绑")
	}

	// 更新卡密信息
	card.DeviceID = nil
	card.Status = 0 // 重置为未使用状态

	result = s.db.Save(&card)
	if result.Error != nil {
		return nil, errors.New("更新卡密信息失败")
	}

	// 返回解绑成功响应
	return &model.UnbindCardResponse{
		Success: true,
		CardNo:  card.CardNo,
		Message: "卡密解绑成功",
	}, nil
}
