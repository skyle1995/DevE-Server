package database

import (
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite" // 使用纯Go实现的SQLite驱动
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// initSQLite 初始化SQLite数据库连接
func initSQLite() (*gorm.DB, error) {
	// 获取数据库文件路径
	dbPath := viper.GetString("database.sqlite.path")

	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, 0755)
		if err != nil {
			log.Errorf("创建数据库目录失败: %v", err)
			return nil, err
		}
	}

	// 配置GORM日志
	gormLogger := logger.New(
		log.StandardLogger(),
		logger.Config{
			SlowThreshold:             time.Second, // 慢SQL阈值
			LogLevel:                  getGormLogLevel(),
			IgnoreRecordNotFoundError: true, // 忽略记录未找到错误
			Colorful:                  true, // 启用彩色日志，错误显示为红色
		},
	)

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Errorf("连接SQLite数据库失败: %v", err)
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("获取数据库连接失败: %v", err)
		return nil, err
	}

	// SQLite连接池设置
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
