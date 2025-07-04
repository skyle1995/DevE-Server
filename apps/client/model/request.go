package model

// ActivateCardRequest 激活卡密请求
type ActivateCardRequest struct {
	CardNo     string                 `json:"card_no" binding:"required"`   // 卡号
	CardKey    string                 `json:"card_key" binding:"required"`  // 卡密
	DeviceID   string                 `json:"device_id" binding:"required"` // 设备ID
	DeviceInfo map[string]interface{} `json:"device_info"`                  // 设备信息
	AppKey     string                 `json:"app_key" binding:"required"`   // 应用密钥
	Timestamp  int64                  `json:"timestamp" binding:"required"` // 时间戳
}

// RebindCardRequest 换绑卡密请求
type RebindCardRequest struct {
	CardNo     string                 `json:"card_no" binding:"required"`   // 卡号
	DeviceID   string                 `json:"device_id" binding:"required"` // 新设备ID
	DeviceInfo map[string]interface{} `json:"device_info"`                  // 设备信息
	AppKey     string                 `json:"app_key" binding:"required"`   // 应用密钥
	Timestamp  int64                  `json:"timestamp" binding:"required"` // 时间戳
}

// UnbindCardRequest 解绑卡密请求
type UnbindCardRequest struct {
	CardNo    string `json:"card_no" binding:"required"`   // 卡号
	AppKey    string `json:"app_key" binding:"required"`   // 应用密钥
	Timestamp int64  `json:"timestamp" binding:"required"` // 时间戳
}

// VerifyAppRequest 验证应用请求
type VerifyAppRequest struct {
	AppKey    string `json:"app_key" binding:"required"`   // 应用密钥
	Timestamp int64  `json:"timestamp" binding:"required"` // 时间戳
}

// HeartbeatRequest 心跳请求
type HeartbeatRequest struct {
	CardNo    string `json:"card_no" binding:"required"`   // 卡号
	DeviceID  string `json:"device_id" binding:"required"` // 设备ID
	AppKey    string `json:"app_key" binding:"required"`   // 应用密钥
	Timestamp int64  `json:"timestamp" binding:"required"` // 时间戳
}
