package model

// GetSettingsRequest 获取系统设置请求参数
type GetSettingsRequest struct {
	Group string `form:"group" json:"group"` // 设置分组，可选
}

// UpdateSettingRequest 更新系统设置请求参数
type UpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`   // 设置键名
	Value string `json:"value" binding:"required"` // 设置值
}

// BatchUpdateSettingsRequest 批量更新系统设置请求参数
type BatchUpdateSettingsRequest struct {
	Settings []UpdateSettingRequest `json:"settings" binding:"required,dive"` // 设置列表
}

// CreateSettingRequest 创建系统设置请求参数
type CreateSettingRequest struct {
	Key         string `json:"key" binding:"required"`    // 设置键名
	Value       string `json:"value" binding:"required"`  // 设置值
	Group       string `json:"group" binding:"required"` // 设置分组
	Description string `json:"description"`               // 设置描述
}
