package model

import "time"

// CreateNoticeRequest 创建通知请求
type CreateNoticeRequest struct {
	Title     string     `json:"title" binding:"required"`    // 通知标题
	Content   string     `json:"content" binding:"required"` // 通知内容
	Level     int        `json:"level" binding:"required"`    // 通知等级：0-普通，1-重要，2-紧急
	Status    int        `json:"status" binding:"required"`   // 状态：0-禁用，1-启用
	StartTime *time.Time `json:"start_time,omitempty"`       // 通知开始时间
	EndTime   *time.Time `json:"end_time,omitempty"`         // 通知结束时间
}

// UpdateNoticeRequest 更新通知请求
type UpdateNoticeRequest struct {
	Title     string     `json:"title,omitempty"`       // 通知标题
	Content   string     `json:"content,omitempty"`     // 通知内容
	Level     *int       `json:"level,omitempty"`       // 通知等级：0-普通，1-重要，2-紧急
	Status    *int       `json:"status,omitempty"`      // 状态：0-禁用，1-启用
	StartTime *time.Time `json:"start_time,omitempty"` // 通知开始时间
	EndTime   *time.Time `json:"end_time,omitempty"`   // 通知结束时间
}

// GetNoticeListRequest 获取通知列表请求
type GetNoticeListRequest struct {
	Page     int    `form:"page" json:"page" binding:"required"`           // 页码
	PageSize int    `form:"page_size" json:"page_size" binding:"required"` // 每页数量
	Title    string `form:"title" json:"title"`                             // 标题，可选，模糊查询
	Level    *int   `form:"level" json:"level"`                             // 等级，可选
	Status   *int   `form:"status" json:"status"`                           // 状态，可选
}

// UpdateNoticeStatusRequest 更新通知状态请求
type UpdateNoticeStatusRequest struct {
	Status int `json:"status" binding:"required"` // 状态：0-禁用，1-启用
}