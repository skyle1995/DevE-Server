package model

import (
	"time"

	"gorm.io/gorm"
)

// Device 设备模型
type Device struct {
	ID          uint                   `gorm:"primaryKey" json:"id"`                                      // 主键ID
	DeviceID    string                 `gorm:"size:100;uniqueIndex;not null" json:"device_id"`              // 设备唯一标识
	DeviceName  string                 `gorm:"size:100" json:"device_name"`                                 // 设备名称
	DeviceType  string                 `gorm:"size:50" json:"device_type"`                                  // 设备类型
	DeviceOS    string                 `gorm:"size:50" json:"device_os"`                                    // 操作系统
	DeviceIP    string                 `gorm:"size:50" json:"device_ip"`                                    // IP地址
	DeviceInfo  map[string]interface{} `gorm:"type:json" json:"device_info"`                                // 设备信息（JSON格式，存储设备详细信息）
	LastActive  time.Time              `json:"last_active"`                                                 // 最后活跃时间
	Status      int                    `gorm:"default:1" json:"status"`                                     // 状态：1-正常, 0-禁用
	AppID       uint                   `json:"app_id"`                                                      // 所属应用ID
	Application App                    `gorm:"foreignKey:AppID" json:"-"`                                   // 所属应用
	CreatedAt   time.Time              `json:"created_at"`                                                  // 创建时间
	UpdatedAt   time.Time              `json:"updated_at"`                                                  // 更新时间
	DeletedAt   gorm.DeletedAt         `gorm:"index" json:"-"`                                             // 删除时间
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}

// UpdateLastActive 更新最后活跃时间
func (d *Device) UpdateLastActive() {
	d.LastActive = time.Now()
}

// IsActive 检查设备是否活跃（最近24小时内有活动）
func (d *Device) IsActive() bool {
	return time.Since(d.LastActive) < 24*time.Hour
}
