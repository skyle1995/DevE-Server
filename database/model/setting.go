package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemSetting 系统设置模型
type SystemSetting struct {
	ID          uint           `gorm:"primaryKey" json:"id"`                    // 主键ID
	Key         string         `gorm:"size:50;uniqueIndex;not null" json:"key"` // 设置键名
	Value       string         `gorm:"size:255" json:"value"`                   // 设置值
	Description string         `gorm:"size:255" json:"description"`             // 设置描述
	Group       string         `gorm:"size:50;default:'system'" json:"group"`   // 设置分组
	CreatedAt   time.Time      `json:"created_at"`                              // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                              // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                          // 删除时间
}

// TableName 指定表名
func (SystemSetting) TableName() string {
	return "settings"
}

// 系统设置默认值
var DefaultSystemSettings = []SystemSetting{
	// 网站基本信息
	{
		Key:         "site_name",
		Value:       "DevE网络验证系统",
		Description: "网站名称",
		Group:       "site",
	},
	{
		Key:         "site_subtitle",
		Value:       "一个基于Golang开发的高性能网络验证管理系统",
		Description: "网站副标题",
		Group:       "site",
	},
	{
		Key:         "site_description",
		Value:       "DevE是一个开发者管理系统，提供应用管理、用户管理、卡密管理等功能，提供软件授权、验证、管理等服务",
		Description: "网站描述",
		Group:       "site",
	},
	{
		Key:         "site_keywords",
		Value:       "DevE,开发者管理系统,应用管理,用户管理,卡密管理",
		Description: "网站关键词",
		Group:       "site",
	},
	{
		Key:         "site_copyright",
		Value:       "© 2025 DevE网络验证系统. 保留所有权利.",
		Description: "网站版权信息",
		Group:       "site",
	},

	// 用户设置
	{
		Key:         "user_enable_login",
		Value:       "1",
		Description: "是否开启普通会员和VIP会员登录(1-开启,0-关闭)",
		Group:       "user",
	},
	{
		Key:         "user_enable_register",
		Value:       "1",
		Description: "是否开放普通会员注册(1-开启,0-关闭)",
		Group:       "user",
	},
	{
		Key:         "user_require_invite_code",
		Value:       "0",
		Description: "注册是否需要邀请码(1-需要,0-不需要)",
		Group:       "user",
	},
	{
		Key:         "user_require_email_verification",
		Value:       "1",
		Description: "注册是否需要邮箱验证(1-需要,0-不需要)",
		Group:       "user",
	},
	{
		Key:         "user_default_status",
		Value:       "1",
		Description: "用户默认状态(1-启用,0-禁用)",
		Group:       "user",
	},
	{
		Key:         "user_remember_me_days",
		Value:       "30",
		Description: "记住密码的有效天数",
		Group:       "user",
	},

	// 邮件设置
	{
		Key:         "mail_enable",
		Value:       "1",
		Description: "是否启用邮件功能(1-启用,0-禁用)",
		Group:       "mail",
	},
	{
		Key:         "mail_smtp_host",
		Value:       "smtp.example.com",
		Description: "SMTP服务器地址",
		Group:       "mail",
	},
	{
		Key:         "mail_smtp_port",
		Value:       "587",
		Description: "SMTP服务器端口",
		Group:       "mail",
	},
	{
		Key:         "mail_smtp_secure",
		Value:       "1",
		Description: "是否使用SSL/TLS(1-使用,0-不使用)",
		Group:       "mail",
	},
	{
		Key:         "mail_from_email",
		Value:       "noreply@example.com",
		Description: "发件人邮箱",
		Group:       "mail",
	},
	{
		Key:         "mail_from_name",
		Value:       "验证系统",
		Description: "发件人名称",
		Group:       "mail",
	},
	{
		Key:         "mail_username",
		Value:       "noreply@example.com",
		Description: "SMTP用户名",
		Group:       "mail",
	},
	{
		Key:         "mail_password",
		Value:       "your_smtp_password",
		Description: "SMTP密码",
		Group:       "mail",
	},

	// 安全设置
	{
		Key:         "security_max_login_attempts",
		Value:       "5",
		Description: "最大登录尝试次数",
		Group:       "security",
	},
	{
		Key:         "security_login_lock_time",
		Value:       "30",
		Description: "登录锁定时间(分钟)",
		Group:       "security",
	},

	// 应用配置
	{
		Key:         "application_key_length",
		Value:       "32",
		Description: "应用密钥长度",
		Group:       "application",
	},
	{
		Key:         "application_max_count",
		Value:       "0",
		Description: "最大应用数量（0表示不限制）",
		Group:       "application",
	},
	{
		Key:         "application_default_status",
		Value:       "1",
		Description: "默认状态（1-启用，0-禁用）",
		Group:       "application",
	},
}
