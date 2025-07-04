package model

import (
	"time"

	dbmodel "github.com/skyle1995/DevE-Server/database/model"
)

// SettingResponse 系统设置响应
type SettingResponse struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	Group       string    `json:"group"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SettingsGroupResponse 按分组的系统设置响应
type SettingsGroupResponse struct {
	Group    string            `json:"group"`
	Settings []SettingResponse `json:"settings"`
}

// FromSetting 将数据库模型转换为响应模型
func FromSetting(setting *dbmodel.SystemSetting) SettingResponse {
	return SettingResponse{
		ID:          setting.ID,
		Key:         setting.Key,
		Value:       setting.Value,
		Description: setting.Description,
		Group:       setting.Group,
		CreatedAt:   setting.CreatedAt,
		UpdatedAt:   setting.UpdatedAt,
	}
}

// FromSettings 将数据库模型列表转换为响应模型列表
func FromSettings(settings []dbmodel.SystemSetting) []SettingResponse {
	responses := make([]SettingResponse, len(settings))
	for i, setting := range settings {
		responses[i] = FromSetting(&setting)
	}
	return responses
}

// GroupSettings 将设置按分组进行分组
func GroupSettings(settings []SettingResponse) []SettingsGroupResponse {
	// 使用map存储分组信息
	groupMap := make(map[string][]SettingResponse)

	// 将设置按分组归类
	for _, setting := range settings {
		groupMap[setting.Group] = append(groupMap[setting.Group], setting)
	}

	// 转换为响应格式
	responses := make([]SettingsGroupResponse, 0, len(groupMap))
	for group, groupSettings := range groupMap {
		responses = append(responses, SettingsGroupResponse{
			Group:    group,
			Settings: groupSettings,
		})
	}

	return responses
}

// SiteInfoResponse 站点信息响应
type SiteInfoResponse struct {
	Name        string `json:"name"`        // 站点名称
	Subtitle    string `json:"subtitle"`    // 站点副标题
	Description string `json:"description"` // 站点描述
	Keywords    string `json:"keywords"`    // 站点关键词
	Copyright   string `json:"copyright"`   // 版权信息
}
