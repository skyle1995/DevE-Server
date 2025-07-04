package model

import "time"

// ActivateCardResponse 激活卡密响应
type ActivateCardResponse struct {
	CardNo        string     `json:"card_no"`        // 卡号
	Status        int        `json:"status"`         // 状态
	Activated     bool       `json:"activated"`      // 是否已激活
	ExpireAt      *time.Time `json:"expire_time"`    // 过期时间
	DeviceBinding bool       `json:"device_binding"` // 是否绑定设备
	BindingInfo   string     `json:"binding_info"`   // 绑定信息
	BindCount     int        `json:"bind_count"`     // 已绑定次数
	MaxBindCount  int        `json:"max_bind_count"` // 最大绑定次数
	CanRebind     bool       `json:"can_rebind"`     // 是否可以换绑
	CanUnbind     bool       `json:"can_unbind"`     // 是否可以解绑
	Message       string     `json:"message"`        // 消息
}

// RebindCardResponse 换绑卡密响应
type RebindCardResponse struct {
	Success  bool   `json:"success"`   // 是否成功
	CardNo   string `json:"card_no"`   // 卡号
	DeviceID string `json:"device_id"` // 设备ID
	Message  string `json:"message"`   // 消息
}

// UnbindCardResponse 解绑卡密响应
type UnbindCardResponse struct {
	Success bool   `json:"success"` // 是否成功
	CardNo  string `json:"card_no"` // 卡号
	Message string `json:"message"` // 消息
}

// VerifyAppResponse 验证应用响应
type VerifyAppResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
}

// VerifyDeviceResponse 验证设备响应
type VerifyDeviceResponse struct {
	Success    bool       `json:"success"`     // 是否成功
	CardNo     string     `json:"card_no"`     // 卡号
	ExpireAt   *time.Time `json:"expire_time"` // 过期时间
	RemainDays int        `json:"remain_days"` // 剩余天数
	Message    string     `json:"message"`     // 消息
}

// HeartbeatResponse 心跳响应
type HeartbeatResponse struct {
	Success    bool       `json:"success"`     // 是否成功
	CardNo     string     `json:"card_no"`     // 卡号
	IsOnline   bool       `json:"is_online"`   // 是否在线
	ExpireAt   *time.Time `json:"expire_time"` // 过期时间
	RemainDays int        `json:"remain_days"` // 剩余天数
	Message    string     `json:"message"`     // 消息
}
