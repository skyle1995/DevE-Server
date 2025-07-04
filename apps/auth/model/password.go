package model

import (
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/crypto"
)

// VerifyPassword 验证密码是否正确
// hashedPassword: 数据库中存储的哈希密码
// plainPassword: 用户输入的明文密码
// 返回值: 密码是否匹配
func VerifyPassword(hashedPassword, plainPassword string) bool {
	return crypto.VerifyPassword(hashedPassword, plainPassword)
}

// GenerateToken 为用户生成JWT令牌
func GenerateToken(userID uint, username string, role int) (string, error) {
	user := &User{
		User: dbmodel.User{
			ID:       userID,
			Username: username,
			Role:     role,
		},
	}
	return user.GenerateToken()
}

// HashPassword 对密码进行哈希处理
// plainPassword: 明文密码
// 返回值: 哈希后的密码和可能的错误
func HashPassword(plainPassword string) (string, error) {
	return crypto.HashPassword(plainPassword)
}
