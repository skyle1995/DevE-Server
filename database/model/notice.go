package model

import (
	"time"

	"gorm.io/gorm"
)

// Notice 系统通知模型
type Notice struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:100;not null" json:"title"`            // 通知标题
	Content   string         `gorm:"type:text;not null" json:"content"`         // 通知内容
	Level     int            `gorm:"default:0" json:"level"`                    // 通知等级：0-普通，1-重要，2-紧急
	Status    int            `gorm:"default:1" json:"status"`                   // 状态：0-禁用，1-启用
	StartTime *time.Time     `gorm:"type:datetime" json:"start_time,omitempty"` // 通知开始时间
	EndTime   *time.Time     `gorm:"type:datetime" json:"end_time,omitempty"`   // 通知结束时间
	UserID    uint           `json:"user_id"`                                   // 发布人ID
	User      User           `gorm:"foreignKey:UserID" json:"-"`                // 发布人
	CreatedAt time.Time      `json:"created_at"`                                // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                            // 删除时间
}

// TableName 指定表名
func (Notice) TableName() string {
	return "notices"
}
