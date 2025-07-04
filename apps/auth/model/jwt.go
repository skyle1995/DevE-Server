package model

import (
	"fmt"
	"time"

	"github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/jwt"
	"github.com/spf13/viper"
)

// User 是auth包中的用户模型，继承自database/model中的User
type User struct {
	model.User
}

// GenerateToken 为用户生成JWT令牌
func (u *User) GenerateToken() (string, error) {
	// 创建JWT配置
	jwtConfig := jwt.DefaultConfig()
	// 从配置中获取JWT密钥
	if viper.IsSet("security.jwt_secret") {
		jwtConfig.SigningKey = viper.GetString("security.jwt_secret")
	}
	// 从配置中获取会话超时时间
	if viper.IsSet("server.session_timeout") {
		jwtConfig.ExpiresTime = time.Duration(viper.GetInt("server.session_timeout")) * time.Hour
	}
	// 创建JWT实例
	jwtInstance := jwt.New(jwtConfig)

	// 将角色转换为字符串
	roleStr := fmt.Sprintf("%d", u.Role)
	// 创建令牌
	tokenString, err := jwtInstance.CreateToken(u.ID, u.Username, roleStr)
	return tokenString, err
}
