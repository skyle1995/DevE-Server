package model

import (
	"time"

	"gorm.io/gorm"
)

// App 应用模型
type App struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`                        // 主键ID
	Name               string         `gorm:"size:100;not null" json:"name"`               // 应用名称
	Description        string         `gorm:"type:text" json:"description"`                // 应用描述
	AppKey             string         `gorm:"size:64;uniqueIndex;not null" json:"app_key"` // 应用密钥
	AppSecret          string         `gorm:"size:64;not null" json:"app_secret"`          // 应用密钥
	Status             int            `gorm:"default:1" json:"status"`                     // 状态：0-禁用，1-启用
	Version            string         `gorm:"size:20" json:"version"`                      // 应用当前版本
	DownloadUrl        string         `gorm:"size:255" json:"download_url"`                // 应用下载地址
	LogoUrl            string         `gorm:"size:255" json:"logo_url"`                    // 应用logo URL
	WebsiteUrl         string         `gorm:"size:255" json:"website_url"`                 // 应用官网地址
	ContactInfo        string         `gorm:"size:255" json:"contact_info"`                // 联系方式
	Notice             string         `gorm:"type:text" json:"notice"`                     // 应用公告内容
	PublicData         string         `gorm:"type:text" json:"public_data"`                // 公有数据，可存储JSON格式的公共配置信息
	PrivateData        string         `gorm:"type:text" json:"private_data"`               // 私有数据，可存储JSON格式的私有配置信息
	DeviceBinding      int            `gorm:"default:0" json:"device_binding"`             // 绑定类型：0-不绑定，1-设备绑定，2-IP绑定
	BillingMode        int            `gorm:"default:0" json:"billing_mode"`               // 计费模式：0-时长计费，1-点数计费
	AllowTrial         int            `gorm:"default:0" json:"allow_trial"`                // 是否允许试用：0-不允许，1-允许
	TrialAmount        int            `gorm:"default:0" json:"trial_amount"`               // 试用额度（时长计费模式下为小时数，点数计费模式下为点数）
	MaxDevices         int            `gorm:"default:0" json:"max_devices"`                // 最大设备数，0表示不限制
	Heartbeat          int            `gorm:"default:0" json:"heartbeat"`                  // 心跳间隔（分钟），0表示不需要心跳
	BindPermission     int            `gorm:"default:0" json:"bind_permission"`            // 绑定权限：0-不允许换绑和解绑，1-允许换绑，2-允许解绑，3-允许换绑和解绑
	UnbindDeductHours  int            `gorm:"default:0" json:"unbind_deduct_hours"`        // 解绑或换绑扣除时长（小时），0表示不扣除
	EncryptionType     int            `gorm:"default:0" json:"encryption_type"`            // 加密类型：0-不加密，1-AES加密，2-RSA加密，3-RC4加密
	EncryptionKey      string         `gorm:"size:255" json:"encryption_key"`              // 加密密钥
	SignatureRequired  int            `gorm:"default:0" json:"signature_required"`         // 是否需要签名：0-不需要，1-需要
	SignatureAlgorithm string         `gorm:"size:20" json:"signature_algorithm"`          // 签名算法：MD5/SHA1/SHA256等
	SignatureKey       string         `gorm:"size:255" json:"signature_key"`               // 签名密钥
	IpWhitelist        string         `gorm:"type:text" json:"ip_whitelist"`               // IP白名单，多个IP用逗号分隔，为空表示不限制
	RequestRateLimit   int            `gorm:"default:0" json:"request_rate_limit"`         // 请求频率限制（次/分钟），0表示不限制
	Timeout            int            `gorm:"default:60" json:"timeout"`                   // 请求超时时间（秒）
	UserID             uint           `json:"user_id"`                                     // 所属用户ID
	User               User           `gorm:"foreignKey:UserID" json:"-"`                  // 所属用户
	CreatedAt          time.Time      `json:"created_at"`                                  // 创建时间
	UpdatedAt          time.Time      `json:"updated_at"`                                  // 更新时间
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`                              // 删除时间
}

// TableName 指定表名
func (App) TableName() string {
	return "apps"
}

// BeforeCreate 创建前的钩子
func (a *App) BeforeCreate(tx *gorm.DB) error {
	// 如果没有设置AppKey和AppSecret，可以在这里自动生成
	// 这里暂时不实现，因为通常会在服务层生成这些值
	return nil
}
