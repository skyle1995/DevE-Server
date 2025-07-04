package database

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// Paginate 分页查询辅助函数
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		if pageSize <= 0 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// GetByID 通过ID获取记录
func GetByID(model interface{}, id uint) error {
	result := GetDB().First(model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("记录不存在: %v", result.Error)
		}
		return fmt.Errorf("查询记录失败: %v", result.Error)
	}
	return nil
}

// Create 创建记录
func Create(model interface{}) error {
	result := GetDB().Create(model)
	if result.Error != nil {
		return fmt.Errorf("创建记录失败: %v", result.Error)
	}
	return nil
}

// Update 更新记录
func Update(model interface{}) error {
	result := GetDB().Save(model)
	if result.Error != nil {
		return fmt.Errorf("更新记录失败: %v", result.Error)
	}
	return nil
}

// Delete 删除记录
func Delete(model interface{}) error {
	result := GetDB().Delete(model)
	if result.Error != nil {
		return fmt.Errorf("删除记录失败: %v", result.Error)
	}
	return nil
}

// Count 获取记录总数
func Count(model interface{}, where ...interface{}) (int64, error) {
	var count int64
	db := GetDB().Model(model)

	if len(where) > 0 {
		db = db.Where(where[0], where[1:]...)
	}

	result := db.Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("获取记录总数失败: %v", result.Error)
	}

	return count, nil
}

// IsExist 检查记录是否存在
func IsExist(model interface{}, where ...interface{}) (bool, error) {
	count, err := Count(model, where...)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetModelName 获取模型名称
func GetModelName(model interface{}) string {
	t := reflect.TypeOf(model)

	// 如果是指针，获取其指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

// Transaction 事务处理
func Transaction(fc func(tx *gorm.DB) error) error {
	return GetDB().Transaction(func(tx *gorm.DB) error {
		return fc(tx)
	})
}
