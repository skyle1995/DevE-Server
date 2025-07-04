package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination 分页结构体
type Pagination struct {
	Page      int         `json:"page" form:"page"`           // 当前页码
	PageSize  int         `json:"page_size" form:"page_size"` // 每页数量
	Total     int64       `json:"total"`                      // 总记录数
	TotalPage int         `json:"total_page"`                 // 总页数
	Data      interface{} `json:"data"`                       // 数据
}

// Config 分页配置
type Config struct {
	DefaultPage     int // 默认页码
	DefaultPageSize int // 默认每页数量
	MaxPageSize     int // 最大每页数量
}

// DefaultConfig 默认分页配置
func DefaultConfig() Config {
	return Config{
		DefaultPage:     1,
		DefaultPageSize: 10,
		MaxPageSize:     100,
	}
}

// New 创建分页实例
func New(page, pageSize int, config ...Config) *Pagination {
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	if page <= 0 {
		page = cfg.DefaultPage
	}

	if pageSize <= 0 {
		pageSize = cfg.DefaultPageSize
	} else if pageSize > cfg.MaxPageSize {
		pageSize = cfg.MaxPageSize
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// NewFromGinContext 从Gin上下文创建分页实例
func NewFromGinContext(c *gin.Context, config ...Config) *Pagination {
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(cfg.DefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(cfg.DefaultPageSize)))

	return New(page, pageSize, cfg)
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *Pagination) GetLimit() int {
	return p.PageSize
}

// SetTotal 设置总记录数并计算总页数
func (p *Pagination) SetTotal(total int64) *Pagination {
	p.Total = total
	p.TotalPage = int(math.Ceil(float64(total) / float64(p.PageSize)))
	return p
}

// SetData 设置数据
func (p *Pagination) SetData(data interface{}) *Pagination {
	p.Data = data
	return p
}

// GetPaginationScope 获取GORM分页Scope
func (p *Pagination) GetPaginationScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}
}

// GetTotalScope 获取GORM总数Scope
func (p *Pagination) GetTotalScope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var total int64
		db.Count(&total)
		p.SetTotal(total)
		return db
	}
}

// Paginate 分页查询辅助函数
func Paginate(c *gin.Context, model interface{}, db *gorm.DB, config ...Config) *Pagination {
	pagination := NewFromGinContext(c, config...)

	// 先获取总数
	var total int64
	db.Model(model).Count(&total)
	pagination.SetTotal(total)

	// 如果总数为0，直接返回空数据
	if total == 0 {
		pagination.SetData([]interface{}{})
		return pagination
	}

	// 查询数据
	var result interface{}
	db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Find(&result)
	pagination.SetData(result)

	return pagination
}

// PaginateWithPreload 带预加载的分页查询辅助函数
func PaginateWithPreload(c *gin.Context, model interface{}, db *gorm.DB, preloads []string, config ...Config) *Pagination {
	pagination := NewFromGinContext(c, config...)

	// 先获取总数
	var total int64
	db.Model(model).Count(&total)
	pagination.SetTotal(total)

	// 如果总数为0，直接返回空数据
	if total == 0 {
		pagination.SetData([]interface{}{})
		return pagination
	}

	// 添加预加载
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	// 查询数据
	var result interface{}
	db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Find(&result)
	pagination.SetData(result)

	return pagination
}

// PaginateWithScopes 带Scopes的分页查询辅助函数
func PaginateWithScopes(c *gin.Context, model interface{}, db *gorm.DB, scopes []func(*gorm.DB) *gorm.DB, config ...Config) *Pagination {
	pagination := NewFromGinContext(c, config...)

	// 应用所有Scopes
	for _, scope := range scopes {
		db = db.Scopes(scope)
	}

	// 先获取总数
	var total int64
	db.Model(model).Count(&total)
	pagination.SetTotal(total)

	// 如果总数为0，直接返回空数据
	if total == 0 {
		pagination.SetData([]interface{}{})
		return pagination
	}

	// 查询数据
	var result interface{}
	db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Find(&result)
	pagination.SetData(result)

	return pagination
}
