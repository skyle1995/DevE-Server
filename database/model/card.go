package model

import (
	"time"

	"gorm.io/gorm"
)

// CardType 卡密类型模型
type CardType struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`                      // 主键ID
	Name                  string         `gorm:"size:50;not null" json:"name"`              // 类型名称
	Description           string         `gorm:"size:255" json:"description"`               // 类型描述
	Duration              int            `gorm:"not null" json:"duration"`                  // 有效时长
	TimeUnit              string         `gorm:"size:10;default:'day'" json:"time_unit"`    // 时间单位（day, month, year）
	ValidDays             int            `gorm:"default:0" json:"valid_days"`               // 有效天数
	Price                 float64        `gorm:"type:decimal(10,2);not null" json:"price"`  // 价格
	Status                int            `gorm:"default:1" json:"status"`                   // 状态：0-禁用，1-启用
	AppID                 uint           `json:"app_id"`                                    // 所属应用ID
	App                   App            `gorm:"foreignKey:AppID" json:"app"`               // 所属应用
	UserID                int            `gorm:"not null" json:"user_id"`                   // 创建者ID
	DefaultMaxRebindCount int            `gorm:"default:0" json:"default_max_rebind_count"` // 默认最大换绑次数
	DefaultMaxUnbindCount int            `gorm:"default:0" json:"default_max_unbind_count"` // 默认最大解绑次数
	MaxBindCount          int            `gorm:"default:0" json:"max_bind_count"`           // 最大换绑/解绑次数
	CreatedAt             time.Time      `json:"created_at"`                                // 创建时间
	UpdatedAt             time.Time      `json:"updated_at"`                                // 更新时间
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`                            // 删除时间
}

// TableName 指定表名
func (CardType) TableName() string {
	return "types"
}

// Card 卡密模型
type Card struct {
	ID             uint           `gorm:"primaryKey" json:"id"`                          // 主键ID
	CardNo         string         `gorm:"size:50;uniqueIndex;not null" json:"card_no"`   // 卡密号
	CardKey        string         `gorm:"size:100;uniqueIndex;not null" json:"card_key"` // 卡密密钥
	TypeID         uint           `json:"type_id"`                                       // 卡密类型ID
	CardType       CardType       `gorm:"foreignKey:TypeID" json:"card_type"`            // 卡密类型
	AppID          uint           `json:"app_id"`                                        // 所属应用ID
	App            App            `gorm:"foreignKey:AppID" json:"app"`                   // 所属应用
	UserID         int            `gorm:"not null" json:"user_id"`                       // 创建者ID
	Status         int            `gorm:"default:0" json:"status"`                       // 状态：0-未使用, 1-已使用, 2-已过期, 3-已禁用
	DeviceID       *string        `json:"device_id"`                                     // 使用设备ID
	Device         *Device        `gorm:"foreignKey:DeviceID" json:"device,omitempty"`   // 使用设备
	BindingInfo    string         `gorm:"type:text" json:"binding_info"`                 // 绑定信息（JSON格式，根据应用的绑定类型存储设备ID或IP地址）
	BindCount      int            `gorm:"default:0" json:"bind_count"`                   // 换绑/解绑次数统计
	MaxRebindCount int            `gorm:"default:0" json:"max_rebind_count"`             // 最大换绑次数
	RebindCount    int            `gorm:"default:0" json:"rebind_count"`                 // 已换绑次数
	MaxUnbindCount int            `gorm:"default:0" json:"max_unbind_count"`             // 最大解绑次数
	UnbindCount    int            `gorm:"default:0" json:"unbind_count"`                 // 已解绑次数
	ActivateAt     *time.Time     `json:"activate_at"`                                   // 激活时间
	ExpireAt       *time.Time     `json:"expire_at"`                                     // 过期时间
	IsOnline       int            `gorm:"default:0" json:"is_online"`                    // 在线状态：0-离线，1-在线
	LastHeartbeat  *time.Time     `json:"last_heartbeat"`                                // 最后心跳时间
	CreatedAt      time.Time      `json:"created_at"`                                    // 创建时间
	UpdatedAt      time.Time      `json:"updated_at"`                                    // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                                // 删除时间
}

// TableName 指定表名
func (Card) TableName() string {
	return "cards"
}

// BeforeCreate 创建前的钩子
func (c *Card) BeforeCreate(tx *gorm.DB) error {
	// 如果没有设置CardNo和CardKey，可以在这里自动生成
	// 这里暂时不实现，因为通常会在服务层生成这些值
	return nil
}

// Activate 激活卡密
func (c *Card) Activate(deviceID string) {
	now := time.Now()
	c.Status = 1
	c.DeviceID = &deviceID
	c.ActivateAt = &now
	c.RebindCount = 0 // 初始化换绑次数

	// 计算过期时间
	if c.CardType.ValidDays > 0 {
		expireAt := now.AddDate(0, 0, c.CardType.ValidDays)
		c.ExpireAt = &expireAt
	}
}

// IsExpired 检查卡密是否已过期
func (c *Card) IsExpired() bool {
	if c.ExpireAt == nil {
		return false
	}
	return time.Now().After(*c.ExpireAt)
}
