package model

import (
	"time"

	dbmodel "github.com/skyle1995/DevE-Server/database/model"
)

// ===== 应用响应模型 =====

// AppResponse 应用响应模型
type AppResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AppKey      string    `json:"app_key"`
	AppSecret   string    `json:"app_secret,omitempty"` // 仅在创建和重新生成密钥时返回
	Status      int       `json:"status"`
	Version     string    `json:"version"`
	DownloadUrl string    `json:"download_url"`
	BillingMode int       `json:"billing_mode"`
	TrialAmount int       `json:"trial_amount"`
	AllowTrial  int       `json:"allow_trial"`
	PublicData  string    `json:"public_data"`  // 公有数据，可存储JSON格式的公共配置信息
	PrivateData string    `json:"private_data"` // 私有数据，可存储JSON格式的私有配置信息
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewAppResponse 从数据库模型创建应用响应模型
func NewAppResponse(app *dbmodel.App) *AppResponse {
	return &AppResponse{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
		AppKey:      app.AppKey,
		AppSecret:   app.AppSecret,
		Status:      app.Status,
		Version:     app.Version,
		DownloadUrl: app.DownloadUrl,
		BillingMode: app.BillingMode,
		TrialAmount: app.TrialAmount,
		AllowTrial:  app.AllowTrial,
		PublicData:  app.PublicData,
		PrivateData: app.PrivateData,
		UserID:      app.UserID,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}

// NewAppResponseWithoutSecret 从数据库模型创建应用响应模型（不包含密钥）
func NewAppResponseWithoutSecret(app *dbmodel.App) *AppResponse {
	response := NewAppResponse(app)
	response.AppSecret = ""
	return response
}

// NewAppResponseList 从数据库模型列表创建应用响应模型列表
func NewAppResponseList(apps []dbmodel.App) []*AppResponse {
	responseList := make([]*AppResponse, len(apps))
	for i, app := range apps {
		responseList[i] = NewAppResponseWithoutSecret(&app)
	}
	return responseList
}
