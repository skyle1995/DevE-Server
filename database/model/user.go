package model

import (
	"time"

	"github.com/skyle1995/DevE-Server/utils/crypto"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                         // 主键ID
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"` // 用户名
	Password  string         `gorm:"size:100;not null" json:"-"`                   // 密码（不返回给前端）
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`            // 邮箱
	Nickname  string         `gorm:"size:50" json:"nickname"`                      // 昵称
	Avatar    string         `gorm:"size:255" json:"avatar"`                       // 头像URL
	Role      int            `gorm:"default:1" json:"role"`                        // 角色：0-管理员, 1-普通会员, 2-VIP会员
	Status    int            `gorm:"default:1" json:"status"`                      // 状态：1-启用, 0-禁用
	LastLogin *time.Time     `json:"last_login"`                                   // 最后登录时间
	CreatedAt time.Time      `json:"created_at"`                                   // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                                   // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                               // 删除时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.hashPassword()
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 只有在密码字段被修改时才重新哈希
	if tx.Statement.Changed("Password") {
		return u.hashPassword()
	}
	return nil
}

// hashPassword 对密码进行哈希处理
func (u *User) hashPassword() error {
	if len(u.Password) == 0 {
		return nil
	}

	hashedPassword, err := crypto.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	return crypto.VerifyPassword(u.Password, password)
}
