package model

import (
	"time"

	dbmodel "github.com/skyle1995/DevE-Server/database/model"
)

// CardTypeResponse 卡密类型响应
type CardTypeResponse struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	Duration              int       `json:"duration"`  // 时长
	TimeUnit              string    `json:"time_unit"` // 时间单位（day, month, year）
	Price                 float64   `json:"price"`
	Status                int       `json:"status"` // 1-启用, 0-禁用
	AppID                 uint      `json:"app_id"`
	AppName               string    `json:"app_name,omitempty"`
	DefaultMaxRebindCount int       `json:"default_max_rebind_count"` // 默认最大换绑次数
	DefaultMaxUnbindCount int       `json:"default_max_unbind_count"` // 默认最大解绑次数
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// CardResponse 卡密响应
type CardResponse struct {
	ID             uint       `json:"id"`
	CardNo         string     `json:"card_no"`
	CardKey        string     `json:"card_key,omitempty"` // 仅在特定情况下返回
	TypeID         uint       `json:"type_id"`
	CardType       string     `json:"card_type"`
	AppID          uint       `json:"app_id"`
	AppName        string     `json:"app_name,omitempty"`
	Status         int        `json:"status"` // 0-未使用, 1-已使用, 2-已过期, 3-已禁用
	DeviceID       *string    `json:"device_id,omitempty"`
	MaxRebindCount int        `json:"max_rebind_count"` // 最大换绑次数
	RebindCount    int        `json:"rebind_count"`     // 已换绑次数
	MaxUnbindCount int        `json:"max_unbind_count"` // 最大解绑次数
	UnbindCount    int        `json:"unbind_count"`     // 已解绑次数
	ActivateAt     *time.Time `json:"activate_at,omitempty"`
	ExpireAt       *time.Time `json:"expire_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// 这里保留给卡密管理相关响应

// GenerateCardResponse 生成卡密响应
type GenerateCardResponse struct {
	Success bool           `json:"success"`
	Count   int            `json:"count"`
	Cards   []CardResponse `json:"cards"`
}

// CardListResponse 卡密列表响应
type CardListResponse struct {
	Total int            `json:"total"`
	Items []CardResponse `json:"items"`
}

// CardTypeListResponse 卡密类型列表响应
type CardTypeListResponse struct {
	Total int                `json:"total"`
	Items []CardTypeResponse `json:"items"`
}

// FromCardType 将数据库卡密类型模型转换为响应模型
func FromCardType(cardType dbmodel.CardType) CardTypeResponse {
	response := CardTypeResponse{
		ID:                    cardType.ID,
		Name:                  cardType.Name,
		Description:           cardType.Description,
		Duration:              cardType.Duration,
		TimeUnit:              cardType.TimeUnit,
		Price:                 cardType.Price,
		Status:                cardType.Status,
		AppID:                 cardType.AppID,
		DefaultMaxRebindCount: cardType.DefaultMaxRebindCount,
		DefaultMaxUnbindCount: cardType.DefaultMaxUnbindCount,
		CreatedAt:             cardType.CreatedAt,
		UpdatedAt:             cardType.UpdatedAt,
	}

	// 如果有应用信息，添加应用名称
	if cardType.App.ID > 0 {
		response.AppName = cardType.App.Name
	}

	return response
}

// FromCard 将数据库卡密模型转换为响应模型
func FromCard(card dbmodel.Card, includeKey bool) CardResponse {
	response := CardResponse{
		ID:             card.ID,
		CardNo:         card.CardNo,
		TypeID:         card.TypeID,
		CardType:       card.CardType.Name,
		AppID:          card.AppID,
		Status:         card.Status,
		MaxRebindCount: card.MaxRebindCount,
		RebindCount:    card.RebindCount,
		MaxUnbindCount: card.MaxUnbindCount,
		UnbindCount:    card.UnbindCount,
		ActivateAt:     card.ActivateAt,
		ExpireAt:       card.ExpireAt,
		CreatedAt:      card.CreatedAt,
	}

	// 仅在特定情况下返回卡密
	if includeKey {
		response.CardKey = card.CardKey
	}

	// 如果有设备信息，添加设备ID
	if card.DeviceID != nil && card.Device != nil {
		deviceID := card.Device.DeviceID
		response.DeviceID = &deviceID
	}

	// 如果有应用信息，添加应用名称
	if card.App.ID > 0 {
		response.AppName = card.App.Name
	}

	return response
}

// FromCards 将数据库卡密模型列表转换为响应模型列表
func FromCards(cards []dbmodel.Card, includeKey bool) []CardResponse {
	responses := make([]CardResponse, len(cards))
	for i, card := range cards {
		responses[i] = FromCard(card, includeKey)
	}
	return responses
}

// FromCardTypes 将数据库卡密类型模型列表转换为响应模型列表
func FromCardTypes(cardTypes []dbmodel.CardType) []CardTypeResponse {
	responses := make([]CardTypeResponse, len(cardTypes))
	for i, cardType := range cardTypes {
		responses[i] = FromCardType(cardType)
	}
	return responses
}
