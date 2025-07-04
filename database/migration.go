package database

import (
	log "github.com/sirupsen/logrus"
	"github.com/skyle1995/DevE-Server/database/model"
	"gorm.io/gorm"
)

// Migration 数据库迁移结构体
type Migration struct {
	db *gorm.DB
}

// NewMigration 创建新的数据库迁移实例
func NewMigration() *Migration {
	return &Migration{db: GetDB()}
}

// AutoMigrate 自动迁移所有模型
func (m *Migration) AutoMigrate() error {
	log.Info("开始自动迁移数据库表结构...")

	// 添加所有需要迁移的模型
	models := []interface{}{
		&model.User{},
		&model.App{},
		&model.CardType{},
		&model.Card{},
		&model.Device{},
		&model.SystemSetting{},
		&model.App{},
		&model.Notice{},
		&model.Logs{},
	}

	for _, model := range models {
		err := m.db.AutoMigrate(model)
		if err != nil {
			log.Errorf("迁移表失败 %T: %v", model, err)
			return err
		}
	}

	log.Info("数据库表结构迁移完成")
	return nil
}

// initAppSettings 初始化应用设置
func (m *Migration) initAppSettings() error {
	// 查询所有应用
	var apps []model.App
	if err := m.db.Find(&apps).Error; err != nil {
		return err
	}

	// 为每个应用创建默认设置
	for _, app := range apps {
		// 检查是否已存在设置
		var count int64
		m.db.Model(&model.App{}).Where("app_id = ?", app.ID).Count(&count)
		if count > 0 {
			continue // 已存在设置，跳过
		}

		// 创建默认设置
		appSetting := model.App{
			ID:             app.ID,
			BillingMode:    0, // 默认使用时长计费
			BindPermission: 3, // 默认允许换绑和解绑
		}

		if err := m.db.Create(&appSetting).Error; err != nil {
			return err
		}
	}

	return nil
}

// InitDefaultData 初始化默认数据
func (m *Migration) InitDefaultData() error {
	// 初始化系统设置
	if err := m.initSystemSettings(); err != nil {
		return err
	}

	// 初始化管理员账户
	if err := m.initAdminUser(); err != nil {
		return err
	}

	// 初始化应用设置
	if err := m.initAppSettings(); err != nil {
		return err
	}

	return nil
}

// initSystemSettings 初始化系统设置
func (m *Migration) initSystemSettings() error {
	// 检查是否已存在系统设置
	var count int64
	m.db.Model(&model.SystemSetting{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 批量创建默认系统设置
	result := m.db.Create(&model.DefaultSystemSettings)
	if result.Error != nil {
		log.Errorf("初始化系统设置失败: %v", result.Error)
		return result.Error
	}

	return nil
}

// initAdminUser 初始化管理员账户
func (m *Migration) initAdminUser() error {
	// 检查是否已存在管理员账户或同名用户
	var count int64
	m.db.Model(&model.User{}).Where("role = ? OR username = ?", 1, "admin").Count(&count)
	if count > 0 {
		return nil
	}

	// 创建默认管理员账户
	admin := model.User{
		Username: "admin",
		Password: "admin123", // 将在BeforeCreate钩子中自动哈希
		Email:    "admin@example.com",
		Nickname: "系统管理员",
		Role:     1,
		Status:   1,
	}

	result := m.db.Create(&admin)
	if result.Error != nil {
		log.Errorf("创建管理员账户失败: %v", result.Error)
		return result.Error
	}

	return nil
}
