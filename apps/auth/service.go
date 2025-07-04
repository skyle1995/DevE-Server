package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/skyle1995/DevE-Server/apps/auth/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"github.com/skyle1995/DevE-Server/utils/jwt"
	"github.com/spf13/viper"
)

// Service 提供认证相关的服务
type Service struct{}

// NewService 创建一个新的认证服务实例
func NewService() *Service {
	return &Service{}
}

// Login 用户登录
func (s *Service) Login(username, password string, remember bool) (dbmodel.User, string, error) {
	// 查询用户
	var dbUser dbmodel.User
	result := database.DB.Where("username = ?", username).First(&dbUser)
	if result.Error != nil {
		return dbmodel.User{}, "", errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if dbUser.Status != 1 {
		return dbmodel.User{}, "", errors.New("账号已被禁用")
	}

	// 验证密码
	if !model.VerifyPassword(dbUser.Password, password) {
		return dbmodel.User{}, "", errors.New("用户名或密码错误")
	}

	// 创建JWT实例
	jwtConfig := jwt.DefaultConfig()
	// 从配置中获取JWT密钥
	if viper.IsSet("security.jwt_secret") {
		jwtConfig.SigningKey = viper.GetString("security.jwt_secret")
	}

	// 设置过期时间
	if remember {
		// 如果记住密码，从数据库获取有效天数设置
		var rememberMeDays int = 30 // 默认30天

		// 查询数据库中的设置
		var setting dbmodel.SystemSetting
		result := database.DB.Where("key = ?", "user_remember_me_days").First(&setting)
		if result.Error == nil && setting.Value != "" {
			// 尝试将设置值转换为整数
			fmt.Sscanf(setting.Value, "%d", &rememberMeDays)
		}

		jwtConfig.ExpiresTime = time.Duration(rememberMeDays) * 24 * time.Hour
	} else {
		// 否则使用配置中的会话超时时间
		if viper.IsSet("server.session_timeout") {
			jwtConfig.ExpiresTime = time.Duration(viper.GetInt("server.session_timeout")) * time.Hour
		} else {
			// 默认24小时
			jwtConfig.ExpiresTime = 24 * time.Hour
		}
	}

	jwtInstance := jwt.New(jwtConfig)

	// 生成JWT令牌
	// 将角色转换为字符串
	roleStr := fmt.Sprintf("%d", dbUser.Role)
	token, err := jwtInstance.CreateToken(dbUser.ID, dbUser.Username, roleStr)
	if err != nil {
		return dbmodel.User{}, "", errors.New("生成令牌失败")
	}

	// 更新用户最后登录时间
	database.DB.Model(&dbUser).Updates(map[string]interface{}{
		"last_login": time.Now(),
	})

	return dbUser, token, nil
}

// Register 处理用户注册
func (s *Service) Register(username, password, email string) error {
	// 检查用户名是否已存在
	var count int64
	database.DB.Model(&dbmodel.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	database.DB.Model(&dbmodel.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已被使用")
	}

	// 检查是否允许注册
	allowRegister := viper.GetBool("system.allow_register")
	if !allowRegister {
		return errors.New("系统当前不允许注册新用户")
	}

	// 哈希密码
	hashedPassword, err := model.HashPassword(password)
	if err != nil {
		return errors.New("密码处理失败")
	}

	// 创建新用户
	newUser := dbmodel.User{
		Username:  username,
		Password:  hashedPassword,
		Email:     email,
		Role:      1, // 默认角色
		Status:    1, // 默认状态：启用
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return errors.New("创建用户失败: " + result.Error.Error())
	}

	return nil
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(userID uint) (map[string]interface{}, error) {
	// 查询用户
	var user dbmodel.User
	result := database.DB.Select("id, username, email, role, status, created_at, last_login").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 构建返回数据
	userInfo := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
		"last_login": user.LastLogin,
	}

	return userInfo, nil
}

// UpdatePassword 更新用户密码
func (s *Service) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	// 查询用户
	var user dbmodel.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !model.VerifyPassword(user.Password, oldPassword) {
		return errors.New("原密码错误")
	}

	// 哈希新密码
	hashedPassword, err := model.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码处理失败")
	}

	// 更新密码
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	result = database.DB.Save(&user)
	if result.Error != nil {
		return errors.New("更新密码失败: " + result.Error.Error())
	}

	return nil
}
