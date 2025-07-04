package database

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// initMySQL 初始化MySQL数据库连接
func initMySQL() (*gorm.DB, error) {
	host := viper.GetString("database.mysql.host")
	port := viper.GetInt("database.mysql.port")
	username := viper.GetString("database.mysql.username")
	password := viper.GetString("database.mysql.password")
	database := viper.GetString("database.mysql.database")
	charset := viper.GetString("database.mysql.charset")
	maxIdleConns := viper.GetInt("database.mysql.max_idle_conns")
	maxOpenConns := viper.GetInt("database.mysql.max_open_conns")

	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, host, port, database, charset)

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Errorf("连接MySQL数据库失败: %v", err)
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("获取数据库连接失败: %v", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// getGormLogLevel 根据应用配置获取GORM日志级别
func getGormLogLevel() logger.LogLevel {
	// 只在错误时打印日志
	return logger.Error
}
