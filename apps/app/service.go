package apps

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/skyle1995/DevE-Server/apps/setting"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/random"
	"gorm.io/gorm"
)

// Service 提供应用相关的服务
type Service struct{}

// NewService 创建一个新的应用服务实例
func NewService() *Service {
	return &Service{}
}

// CreateApp 创建新应用
func (s *Service) CreateApp(
	userID uint,
	name string,
	description string,
	version string,
	downloadUrl string,
	billingMode int,
	trialAmount int,
	allowTrial int,
	publicData string,
	privateData string,
) (*dbmodel.App, error) {
	// 检查用户是否存在
	var user dbmodel.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户创建的应用数量是否超过限制
	var count int64
	database.DB.Model(&dbmodel.App{}).Where("user_id = ?", userID).Count(&count)

	// 创建系统设置服务实例
	settingService := setting.NewService()

	// 从系统设置中获取最大应用数量
	maxApps := int64(0) // 默认不限制
	maxAppsSetting, maxErr := settingService.GetSettingByKey("application_max_count")
	if maxErr == nil {
		// 转换设置值为整数
		if ma, parseErr := strconv.ParseInt(maxAppsSetting.Value, 10, 64); parseErr == nil {
			maxApps = ma
		}
	}

	if maxApps > 0 && count >= maxApps {
		return nil, fmt.Errorf("您最多只能创建%d个应用", maxApps)
	}

	// 从系统设置中获取应用密钥长度
	keyLength := 32 // 默认长度
	keySetting, keyErr := settingService.GetSettingByKey("application_key_length")
	if keyErr == nil {
		// 转换设置值为整数
		if kl, parseErr := strconv.Atoi(keySetting.Value); parseErr == nil && kl > 0 {
			keyLength = kl
		}
	}

	appKey := random.GenerateRandomString(keyLength)
	appSecret := random.GenerateRandomString(keyLength)

	// 从系统设置中获取应用默认状态
	defaultStatus := 1 // 默认启用
	statusSetting, err := settingService.GetSettingByKey("application_default_status")
	if err == nil {
		// 转换设置值为整数
		if ds, parseErr := strconv.Atoi(statusSetting.Value); parseErr == nil && (ds == 0 || ds == 1) {
			defaultStatus = ds
		}
	}

	// 创建新应用
	newApp := dbmodel.App{
		Name:        name,
		Description: description,
		AppKey:      appKey,
		AppSecret:   appSecret,
		Status:      defaultStatus, // 使用系统设置的默认状态
		Version:     version,
		DownloadUrl: downloadUrl,
		BillingMode: billingMode,
		TrialAmount: trialAmount,
		AllowTrial:  allowTrial,
		PublicData:  publicData,
		PrivateData: privateData,
		UserID:      userID,
	}

	result = database.DB.Create(&newApp)
	if result.Error != nil {
		return nil, errors.New("创建应用失败: " + result.Error.Error())
	}

	return &newApp, nil
}

// GetAppList 获取用户的应用列表
func (s *Service) GetAppList(userID uint, page, pageSize int) ([]dbmodel.App, int64, error) {
	var apps []dbmodel.App
	var total int64

	// 计算总数
	database.DB.Model(&dbmodel.App{}).Where("user_id = ?", userID).Count(&total)

	// 获取分页数据
	offset := (page - 1) * pageSize
	result := database.DB.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&apps)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, 0, result.Error
	}

	return apps, total, nil
}

// GetAppByID 获取应用详情
func (s *Service) GetAppByID(userID, appID uint) (*dbmodel.App, error) {
	var app dbmodel.App

	// 查询应用
	result := database.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("应用不存在或无权访问")
		}
		return nil, result.Error
	}

	return &app, nil
}

// UpdateApp 更新应用信息
func (s *Service) UpdateApp(
	userID uint,
	appID uint,
	name string,
	description string,
	version string,
	downloadUrl string,
	billingMode int,
	trialAmount int,
	allowTrial int,
	status int,
	publicData string,
	privateData string,
) (*dbmodel.App, error) {
	// 查询应用是否存在且属于该用户
	var app dbmodel.App
	result := database.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("应用不存在或无权访问")
		}
		return nil, result.Error
	}

	// 更新应用信息
	updates := map[string]interface{}{}

	if name != "" {
		updates["name"] = name
	}

	if description != "" {
		updates["description"] = description
	}

	if version != "" {
		updates["version"] = version
	}

	if downloadUrl != "" {
		updates["download_url"] = downloadUrl
	}

	if billingMode >= 0 && billingMode <= 1 {
		updates["billing_mode"] = billingMode
	}

	if trialAmount >= 0 {
		updates["trial_amount"] = trialAmount
	}

	if allowTrial >= 0 {
		updates["allow_trial"] = allowTrial
	}

	// 只允许在0和1之间切换状态
	if status == 0 || status == 1 {
		updates["status"] = status
	}

	// 更新公有数据和私有数据
	if publicData != "" {
		updates["public_data"] = publicData
	}

	if privateData != "" {
		updates["private_data"] = privateData
	}

	if len(updates) > 0 {
		result = database.DB.Model(&app).Updates(updates)
		if result.Error != nil {
			return nil, errors.New("更新应用失败: " + result.Error.Error())
		}

		// 重新获取更新后的应用信息
		database.DB.First(&app, app.ID)
	}

	return &app, nil
}

// DeleteApp 删除应用
func (s *Service) DeleteApp(userID, appID uint) error {
	// 查询应用是否存在且属于该用户
	var app dbmodel.App
	result := database.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("应用不存在或无权访问")
		}
		return result.Error
	}

	// 删除应用（软删除）
	result = database.DB.Delete(&app)
	if result.Error != nil {
		return errors.New("删除应用失败: " + result.Error.Error())
	}

	return nil
}

// RegenerateAppSecret 重新生成应用密钥
func (s *Service) RegenerateAppSecret(userID, appID uint) (*dbmodel.App, error) {
	// 查询应用是否存在且属于该用户
	var app dbmodel.App
	result := database.DB.Where("id = ? AND user_id = ?", appID, userID).First(&app)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("应用不存在或无权访问")
		}
		return nil, result.Error
	}

	// 从系统设置中获取应用密钥长度并生成新的应用密钥
	keyLength := 32 // 默认长度
	settingService := setting.NewService()
	keySetting, keyErr := settingService.GetSettingByKey("application_key_length")
	if keyErr == nil {
		// 转换设置值为整数
		if kl, parseErr := strconv.Atoi(keySetting.Value); parseErr == nil && kl > 0 {
			keyLength = kl
		}
	}

	newAppSecret := random.GenerateRandomString(keyLength)

	// 更新应用密钥
	result = database.DB.Model(&app).Update("app_secret", newAppSecret)
	if result.Error != nil {
		return nil, errors.New("更新应用密钥失败: " + result.Error.Error())
	}

	// 重新获取更新后的应用信息
	database.DB.First(&app, app.ID)

	return &app, nil
}
