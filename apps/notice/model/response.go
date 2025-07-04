package model

import (
	"time"

	dbmodel "github.com/skyle1995/DevE-Server/database/model"
)

// NoticeResponse 通知响应
type NoticeResponse struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`      // 通知标题
	Content   string     `json:"content"`    // 通知内容
	Level     int        `json:"level"`      // 通知等级：0-普通，1-重要，2-紧急
	Status    int        `json:"status"`     // 状态：0-禁用，1-启用
	StartTime *time.Time `json:"start_time"` // 通知开始时间
	EndTime   *time.Time `json:"end_time"`   // 通知结束时间
	UserID    uint       `json:"user_id"`    // 发布人ID
	UserName  string     `json:"user_name"`  // 发布人用户名
	CreatedAt time.Time  `json:"created_at"` // 创建时间
	UpdatedAt time.Time  `json:"updated_at"` // 更新时间
}

// NoticeListResponse 通知列表响应
type NoticeListResponse struct {
	Total int              `json:"total"` // 总数
	List  []NoticeResponse `json:"list"`  // 列表
}

// ConvertFromDBModel 从数据库模型转换为响应模型
func ConvertFromDBModel(notice dbmodel.Notice, userName string) NoticeResponse {
	return NoticeResponse{
		ID:        notice.ID,
		Title:     notice.Title,
		Content:   notice.Content,
		Level:     notice.Level,
		Status:    notice.Status,
		StartTime: notice.StartTime,
		EndTime:   notice.EndTime,
		UserID:    notice.UserID,
		UserName:  userName,
		CreatedAt: notice.CreatedAt,
		UpdatedAt: notice.UpdatedAt,
	}
}
