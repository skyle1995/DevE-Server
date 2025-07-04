package user

import (
	"errors"
	"strconv"

	authmodel "github.com/skyle1995/DevE-Server/apps/auth/model"
	usermodel "github.com/skyle1995/DevE-Server/apps/user/model"
	"github.com/skyle1995/DevE-Server/database"
	dbmodel "github.com/skyle1995/DevE-Server/database/model"
	"gorm.io/gorm"
)

// Service 提供用户相关的服务
type Service struct{}

// NewService 创建一个新的用户服务实例
func NewService() *Service {
	return &Service{}
}

// GetUserList 获取用户列表
func (s *Service) GetUserList(roleStr, statusStr, pageStr, limitStr string) ([]dbmodel.User, int64, error) {
	// 构建查询
	query := database.DB.Model(&dbmodel.User{})

	// 应用过滤条件
	if roleStr != "" {
		role, err := strconv.Atoi(roleStr)
		if err == nil {
			query = query.Where("role = ?", role)
		}
	}

	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			query = query.Where("status = ?", status)
		}
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// 查询用户列表
	var users []dbmodel.User
	result := query.Select("id, username, email, role, status, created_at, last_login").Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// GetUserByID 根据ID获取用户信息
func (s *Service) GetUserByID(userIDStr string) (*dbmodel.User, error) {
	// 转换用户ID
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}

	// 查询用户
	var user dbmodel.User
	result := database.DB.Select("id, username, email, role, status, created_at, last_login").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}

	return &user, nil
}

// CreateUser 创建新用户
func (s *Service) CreateUser(username, password, email, roleStr, statusStr string) error {
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

	// 转换角色和状态
	role, err := strconv.Atoi(roleStr)
	if err != nil {
		return errors.New("无效的角色值")
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil {
		return errors.New("无效的状态值")
	}

	// 哈希密码
	hashedPassword, err := authmodel.HashPassword(password)
	if err != nil {
		return errors.New("密码处理失败")
	}

	// 创建新用户
	newUser := dbmodel.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Role:     role,
		Status:   status,
	}

	// 保存到数据库
	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return errors.New("创建用户失败: " + result.Error.Error())
	}

	return nil
}

// UpdateUser 更新用户信息
func (s *Service) UpdateUser(userIDStr, email, roleStr, statusStr string) error {
	// 转换用户ID
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return errors.New("无效的用户ID")
	}

	// 查询用户
	var user dbmodel.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return result.Error
	}

	// 准备更新数据
	updates := make(map[string]interface{})

	// 更新邮箱（如果提供）
	if email != "" {
		// 检查邮箱是否已被其他用户使用
		var count int64
		database.DB.Model(&dbmodel.User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
		if count > 0 {
			return errors.New("邮箱已被其他用户使用")
		}
		updates["email"] = email
	}

	// 更新角色（如果提供）
	if roleStr != "" {
		role, err := strconv.Atoi(roleStr)
		if err != nil {
			return errors.New("无效的角色值")
		}
		updates["role"] = role
	}

	// 更新状态（如果提供）
	if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			return errors.New("无效的状态值")
		}
		updates["status"] = status
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 更新用户
	result = database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return errors.New("更新用户失败: " + result.Error.Error())
	}

	return nil
}

// DeleteUser 删除用户
func (s *Service) DeleteUser(userIDStr string) error {
	// 转换用户ID
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return errors.New("无效的用户ID")
	}

	// 查询用户
	var user dbmodel.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return result.Error
	}

	// 删除用户
	result = database.DB.Delete(&user)
	if result.Error != nil {
		return errors.New("删除用户失败: " + result.Error.Error())
	}

	return nil
}

// GetUserProfile 获取用户个人资料
func (s *Service) GetUserProfile(userID uint) (map[string]interface{}, error) {
	// 查询用户
	var user dbmodel.User
	result := database.DB.Select("id, username, email, role, status, created_at, last_login").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}

	// 构建返回数据
	profile := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"role":       user.Role,
		"status":     user.Status,
		"created_at": user.CreatedAt,
		"last_login": user.LastLogin,
	}

	return profile, nil
}

// UpdateUserProfile 更新用户个人资料
func (s *Service) UpdateUserProfile(userID uint, email string) error {
	// 查询用户
	var user dbmodel.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return result.Error
	}

	// 准备更新数据
	updates := make(map[string]interface{})

	// 更新邮箱（如果提供）
	if email != "" {
		// 检查邮箱是否已被其他用户使用
		var count int64
		database.DB.Model(&dbmodel.User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
		if count > 0 {
			return errors.New("邮箱已被其他用户使用")
		}
		updates["email"] = email
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 更新用户
	result = database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return errors.New("更新个人资料失败: " + result.Error.Error())
	}

	return nil
}

// GetUserPermission 获取用户权限
func (s *Service) GetUserPermission(userID uint) (*usermodel.UserPermissionResponse, error) {
	// 查询用户
	var user dbmodel.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}

	// 创建权限响应
	permission := &usermodel.UserPermissionResponse{
		Role: user.Role,
	}

	// 设置角色名称
	switch user.Role {
	case 0:
		permission.RoleName = "管理员"
	case 1:
		permission.RoleName = "普通会员"
	case 2:
		permission.RoleName = "VIP会员"
	default:
		permission.RoleName = "未知角色"
	}

	// 设置权限
	permission.Permissions.CanManageUsers = user.Role == 0                         // 只有管理员可以管理用户
	permission.Permissions.CanAccessVipFeatures = user.Role == 0 || user.Role == 2 // 管理员和VIP会员可以访问VIP功能

	return permission, nil
}
