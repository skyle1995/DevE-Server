package crypto

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行哈希处理
// plainPassword: 明文密码
// 返回值: 哈希后的密码字符串和可能的错误
func HashPassword(plainPassword string) (string, error) {
	// 从配置中获取bcrypt成本，如果未配置则使用默认值
	cost := viper.GetInt("security.bcrypt_cost")
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword 验证密码是否正确
// hashedPassword: 数据库中存储的哈希密码
// plainPassword: 用户输入的明文密码
// 返回值: 密码是否匹配
func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}