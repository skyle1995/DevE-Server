package model

// 这里保留给卡密管理相关请求

// 用户卡密类型管理相关请求

// CreateCardTypeRequest 创建卡密类型请求
type CreateCardTypeRequest struct {
	Name                  string `json:"name" binding:"required"`                     // 卡密类型名称
	Duration              int    `json:"duration" binding:"required"`                 // 时长
	TimeUnit              string `json:"time_unit" binding:"required"`                // 时间单位（day, month, year）
	AppID                 int    `json:"app_id" binding:"required"`                   // 所属应用ID
	Status                int    `json:"status" binding:"required"`                   // 状态：0-禁用，1-启用
	DefaultMaxRebindCount int    `json:"default_max_rebind_count" binding:"required"` // 默认最大换绑次数
	DefaultMaxUnbindCount int    `json:"default_max_unbind_count" binding:"required"` // 默认最大解绑次数
}

// GetCardTypeListRequest 获取卡密类型列表请求
type GetCardTypeListRequest struct {
	Page     int    `form:"page" json:"page" binding:"required"`           // 页码
	PageSize int    `form:"page_size" json:"page_size" binding:"required"` // 每页数量
	AppID    int    `form:"app_id" json:"app_id"`                          // 应用ID，可选
	Status   int    `form:"status" json:"status"`                          // 状态，可选
	Name     string `form:"name" json:"name"`                              // 名称，可选，模糊查询
}

// UpdateCardTypeRequest 更新卡密类型请求
type UpdateCardTypeRequest struct {
	ID                    int    `json:"id" binding:"required"`    // 卡密类型ID
	Name                  string `json:"name"`                     // 卡密类型名称
	Duration              int    `json:"duration"`                 // 时长
	TimeUnit              string `json:"time_unit"`                // 时间单位（day, month, year）
	AppID                 int    `json:"app_id"`                   // 所属应用ID
	Status                int    `json:"status"`                   // 状态：0-禁用，1-启用
	DefaultMaxRebindCount int    `json:"default_max_rebind_count"` // 默认最大换绑次数
	DefaultMaxUnbindCount int    `json:"default_max_unbind_count"` // 默认最大解绑次数
}

// DeleteCardTypeRequest 删除卡密类型请求
type DeleteCardTypeRequest struct {
	ID int `json:"id" binding:"required"` // 卡密类型ID
}

// 用户卡密管理相关请求

// GenerateCardRequest 生成卡密请求
type GenerateCardRequest struct {
	TypeID         int    `json:"type_id" binding:"required"` // 卡密类型ID
	Count          int    `json:"count" binding:"required"`   // 生成数量
	Prefix         string `json:"prefix"`                     // 卡号前缀，可选
	MaxRebindCount *int   `json:"max_rebind_count"`           // 最大换绑次数，可选，不传则使用卡密类型默认值
	MaxUnbindCount *int   `json:"max_unbind_count"`           // 最大解绑次数，可选，不传则使用卡密类型默认值
	KeyLength      *int   `json:"key_length"`                 // 卡密长度，可选，不传则使用默认值16
}

// GetCardListRequest 获取卡密列表请求
type GetCardListRequest struct {
	Page     int    `form:"page" json:"page" binding:"required"`           // 页码
	PageSize int    `form:"page_size" json:"page_size" binding:"required"` // 每页数量
	AppID    int    `form:"app_id" json:"app_id"`                          // 应用ID，可选
	TypeID   int    `form:"type_id" json:"type_id"`                        // 卡密类型ID，可选
	Status   int    `form:"status" json:"status"`                          // 状态，可选
	CardNo   string `form:"card_no" json:"card_no"`                        // 卡号，可选，模糊查询
}

// UpdateCardRequest 更新卡密请求
type UpdateCardRequest struct {
	ID             int  `json:"id" binding:"required"` // 卡密ID
	TypeID         int  `json:"type_id"`               // 卡密类型ID，可选
	Status         int  `json:"status"`                // 状态，可选
	MaxRebindCount *int `json:"max_rebind_count"`      // 最大换绑次数，可选
	MaxUnbindCount *int `json:"max_unbind_count"`      // 最大解绑次数，可选
}

// DeleteCardRequest 删除卡密请求
type DeleteCardRequest struct {
	ID int `json:"id" binding:"required"` // 卡密ID
}

// ExportCardRequest 导出卡密请求
type ExportCardRequest struct {
	TypeID int  `json:"type_id"` // 卡密类型ID，可选
	Used   bool `json:"used"`    // 是否已使用，可选
}
