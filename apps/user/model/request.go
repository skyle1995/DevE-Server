package model

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Role     int    `json:"role" binding:"required,oneof=0 1"`
	Status   int    `json:"status" binding:"required,oneof=0 1"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Role   int    `json:"role" binding:"omitempty,oneof=0 1"`
	Status int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// UpdateUserProfileRequest 更新用户个人资料请求
type UpdateUserProfileRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

// UserListQuery 用户列表查询参数
type UserListQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	Role     int `form:"role" binding:"omitempty,oneof=0 1"`
	Status   int `form:"status" binding:"omitempty,oneof=0 1"`
}