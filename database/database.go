package database

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/skyle1995/DevE-Server/database/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 全局数据库连接实例
var (
	DB     *gorm.DB
	dbOnce sync.Once
)

// Init 初始化数据库连接
func Init() {
	dbType := viper.GetString("database.type")

	dbOnce.Do(func() {
		var err error
		switch dbType {
		case "MySQL":
			DB, err = initMySQL()
		case "SQLite":
			DB, err = initSQLite()
		default:
			log.Warnf("未知的数据库类型: %s，默认使用SQLite", dbType)
			DB, err = initSQLite()
		}

		if err != nil {
			log.Fatalf("数据库初始化失败: %v", err)
		}

		// 初始化数据库表结构
		initTables()
	})
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	if DB == nil {
		log.Warn("数据库连接尚未初始化，正在尝试初始化...")
		Init()
	}
	return DB
}

// Close 关闭数据库连接
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Errorf("获取数据库连接失败: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Errorf("关闭数据库连接失败: %v", err)
			return
		}
	}
}

// initDefaultSystemSettings 初始化默认系统设置
func initDefaultSystemSettings() {
	// 检查是否已存在系统设置
	var count int64
	DB.Model(&model.SystemSetting{}).Count(&count)
	if count > 0 {
		return
	}

	// 批量创建默认系统设置
	result := DB.Create(&model.DefaultSystemSettings)
	if result.Error != nil {
		log.Errorf("初始化默认系统设置失败: %v", result.Error)
		return
	}
}

// initTables 初始化数据库表结构
func initTables() {
	// 自动迁移表结构
	models := []interface{}{
		&model.User{},
		&model.App{},
		&model.CardType{},
		&model.Card{},
		&model.Device{},
		&model.SystemSetting{},
	}

	for _, model := range models {
		err := DB.AutoMigrate(model)
		if err != nil {
			log.Errorf("自动迁移表结构失败: %v", err)
		}
	}

	// 初始化默认系统设置
	initDefaultSystemSettings()
}
