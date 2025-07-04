package setting

import (
	"errors"

	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"gorm.io/gorm"
)

// Service 系统设置服务
type Service struct{}

// NewService 创建一个新的系统设置服务实例
func NewService() *Service {
	return &Service{}
}

// GetSettings 获取系统设置列表
func (s *Service) GetSettings(group string) ([]dbmodel.SystemSetting, error) {
	var settings []dbmodel.SystemSetting
	query := database.DB

	// 如果指定了分组，则按分组筛选
	if group != "" {
		query = query.Where("group = ?", group)
	}

	// 查询设置
	result := query.Find(&settings)
	if result.Error != nil {
		return nil, result.Error
	}

	return settings, nil
}

// GetSettingByKey 根据键名获取系统设置
func (s *Service) GetSettingByKey(key string) (*dbmodel.SystemSetting, error) {
	var setting dbmodel.SystemSetting
	result := database.DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("设置不存在")
		}
		return nil, result.Error
	}

	return &setting, nil
}

// UpdateSetting 更新系统设置
func (s *Service) UpdateSetting(key, value string) error {
	// 检查设置是否存在
	var setting dbmodel.SystemSetting
	result := database.DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("设置不存在")
		}
		return result.Error
	}

	// 更新设置值
	setting.Value = value
	result = database.DB.Save(&setting)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CreateSetting 创建系统设置
func (s *Service) CreateSetting(key, value, group, description string) error {
	// 检查键名是否已存在
	var count int64
	database.DB.Model(&dbmodel.SystemSetting{}).Where("key = ?", key).Count(&count)
	if count > 0 {
		return errors.New("键名已存在")
	}

	// 创建新设置
	newSetting := dbmodel.SystemSetting{
		Key:         key,
		Value:       value,
		Group:       group,
		Description: description,
	}

	// 保存到数据库
	result := database.DB.Create(&newSetting)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteSetting 删除系统设置
func (s *Service) DeleteSetting(key string) error {
	// 检查设置是否存在
	var setting dbmodel.SystemSetting
	result := database.DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("设置不存在")
		}
		return result.Error
	}

	// 删除设置
	result = database.DB.Delete(&setting)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetSiteInfo 获取站点信息
func (s *Service) GetSiteInfo() (map[string]string, error) {
	// 获取site分组的所有设置
	settings, err := s.GetSettings("site")
	if err != nil {
		return nil, err
	}

	// 将设置转换为map
	siteInfo := make(map[string]string)
	for _, setting := range settings {
		// 去掉键名中的site_前缀
		key := setting.Key
		if len(key) > 5 && key[:5] == "site_" {
			key = key[5:]
		}
		siteInfo[key] = setting.Value
	}

	return siteInfo, nil
}
