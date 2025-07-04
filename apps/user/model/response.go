package model

import (
	"time"

	"github.com/skyle1995/DevE-Server/database/model"
)

// UserResponse 用户信息响应
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      int        `json:"role"`
	Status    int        `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64          `json:"total"`
	List  []UserResponse `json:"list"`
}

// FromUser 从数据库用户模型转换为响应模型
func FromUser(user *model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		LastLogin: user.LastLogin,
	}
}

// FromUsers 从数据库用户模型列表转换为响应模型列表
func FromUsers(users []model.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = FromUser(&user)
	}
	return responses
}

// UserPermissionResponse 用户权限响应
type UserPermissionResponse struct {
	Role        int    `json:"role"`      // 角色：0-管理员, 1-普通会员, 2-VIP会员
	RoleName    string `json:"role_name"` // 角色名称
	Permissions struct {
		CanManageUsers       bool `json:"can_manage_users"`        // 是否可以管理用户
		CanAccessVipFeatures bool `json:"can_access_vip_features"` // 是否可以访问VIP功能
	} `json:"permissions"`
}
