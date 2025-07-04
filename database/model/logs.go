package model

import (
	"time"

	"gorm.io/gorm"
)

// Log 日志模型
type Logs struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Type         int            `json:"type"`          // 日志类型：1-登录日志, 2-操作日志，3-系统日志，4-应用日志
	AppID        uint           `json:"app_id"`        // 应用ID，仅应用日志有效
	UserID       uint           `json:"user_id"`       // 用户ID，仅操作日志有效
	Method       string         `json:"method"`        // 请求方法
	Path         string         `json:"path"`          // 请求路径
	IP           string         `json:"ip"`            // 请求IP
	UserAgent    string         `json:"user_agent"`    // 用户代理
	StatusCode   int            `json:"status_code"`   // 状态码
	Latency      int64          `json:"latency"`       // 请求耗时（毫秒）
	RequestBody  string         `json:"request_body"`  // 请求体
	ResponseBody string         `json:"response_body"` // 响应体
	Content      string         `json:"content"`       // 日志内容，主要用于系统日志
}

// TableName 设置表名
func (Logs) TableName() string {
	return "logs"
}

// 日志类型常量
const (
	LogTypeLogin     = 1 // 登录日志
	LogTypeOperation = 2 // 操作日志
	LogTypeSystem    = 3 // 系统日志
	LogTypeApp       = 4 // 应用日志
)
