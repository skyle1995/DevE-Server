package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger 日志记录器结构体
type Logger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

// LogLevel 日志级别类型
type LogLevel string

// 日志级别常量
const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
)

// LogFormat 日志格式类型
type LogFormat string

// 日志格式常量
const (
	TextFormat LogFormat = "text"
	JSONFormat LogFormat = "json"
)

// Config 日志配置结构体
type Config struct {
	Level        LogLevel
	Format       LogFormat
	FilePath     string
	MaxSize      int  // 单个日志文件最大大小（MB）
	MaxBackups   int  // 最大保留的旧日志文件数量
	MaxAge       int  // 日志文件保留的最大天数
	Compress     bool // 是否压缩旧日志文件
	ReportCaller bool // 是否记录调用者信息
}

// DefaultConfig 默认日志配置
func DefaultConfig() Config {
	return Config{
		Level:        InfoLevel,
		Format:       TextFormat,
		FilePath:     "",
		MaxSize:      100,
		MaxBackups:   3,
		MaxAge:       28,
		Compress:     false,
		ReportCaller: true,
	}
}

// New 创建新的日志记录器
func New(config Config) (*Logger, error) {
	logrusLogger := logrus.New()

	// 设置日志级别
	level, err := parseLevel(config.Level)
	if err != nil {
		return nil, err
	}
	logrusLogger.SetLevel(level)

	// 设置日志格式
	switch config.Format {
	case TextFormat:
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	case JSONFormat:
		logrusLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}

	// 设置输出
	if config.FilePath != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		// 打开日志文件
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}

		// 同时输出到文件和标准输出
		logrusLogger.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		// 只输出到标准输出
		logrusLogger.SetOutput(os.Stdout)
	}

	// 设置是否记录调用者信息
	logrusLogger.SetReportCaller(config.ReportCaller)

	return &Logger{
		logger: logrusLogger,
		fields: logrus.Fields{},
	}, nil
}

// parseLevel 解析日志级别
func parseLevel(level LogLevel) (logrus.Level, error) {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel, nil
	case InfoLevel:
		return logrus.InfoLevel, nil
	case WarnLevel:
		return logrus.WarnLevel, nil
	case ErrorLevel:
		return logrus.ErrorLevel, nil
	case FatalLevel:
		return logrus.FatalLevel, nil
	case PanicLevel:
		return logrus.PanicLevel, nil
	default:
		return logrus.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// WithField 添加单个字段
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &Logger{
		logger: l.logger,
		fields: newFields,
	}
}

// WithFields 添加多个字段
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	return &Logger{
		logger: l.logger,
		fields: newFields,
	}
}

// Debug 记录调试级别日志
func (l *Logger) Debug(args ...interface{}) {
	l.logger.WithFields(l.fields).Debug(args...)
}

// Debugf 记录格式化的调试级别日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Debugf(format, args...)
}

// Info 记录信息级别日志
func (l *Logger) Info(args ...interface{}) {
	l.logger.WithFields(l.fields).Info(args...)
}

// Infof 记录格式化的信息级别日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Infof(format, args...)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(args ...interface{}) {
	l.logger.WithFields(l.fields).Warn(args...)
}

// Warnf 记录格式化的警告级别日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Warnf(format, args...)
}

// Error 记录错误级别日志
func (l *Logger) Error(args ...interface{}) {
	l.logger.WithFields(l.fields).Error(args...)
}

// Errorf 记录格式化的错误级别日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Errorf(format, args...)
}

// Fatal 记录致命级别日志
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.WithFields(l.fields).Fatal(args...)
}

// Fatalf 记录格式化的致命级别日志
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Fatalf(format, args...)
}

// Panic 记录恐慌级别日志
func (l *Logger) Panic(args ...interface{}) {
	l.logger.WithFields(l.fields).Panic(args...)
}

// Panicf 记录格式化的恐慌级别日志
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Panicf(format, args...)
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) error {
	logrusLevel, err := parseLevel(level)
	if err != nil {
		return err
	}
	l.logger.SetLevel(logrusLevel)
	return nil
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() LogLevel {
	switch l.logger.GetLevel() {
	case logrus.DebugLevel:
		return DebugLevel
	case logrus.InfoLevel:
		return InfoLevel
	case logrus.WarnLevel:
		return WarnLevel
	case logrus.ErrorLevel:
		return ErrorLevel
	case logrus.FatalLevel:
		return FatalLevel
	case logrus.PanicLevel:
		return PanicLevel
	default:
		return InfoLevel
	}
}

// SetOutput 设置日志输出
func (l *Logger) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

// SetFormatter 设置日志格式化器
func (l *Logger) SetFormatter(format LogFormat) {
	switch format {
	case TextFormat:
		l.logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	case JSONFormat:
		l.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}

// GetCallerInfo 获取调用者信息
func GetCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}

	shortFile := file
	for _, dir := range []string{"github.com", "golang.org", "gopkg.in"} {
		if i := strings.Index(file, dir); i >= 0 {
			shortFile = file[i:]
			break
		}
	}

	return fmt.Sprintf("%s:%d", shortFile, line)
}

// 全局日志实例
var defaultLogger *Logger

// 初始化默认日志记录器
func init() {
	config := DefaultConfig()
	logger, err := New(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize default logger: %v", err))
	}
	defaultLogger = logger
}

// Default 获取默认日志记录器
func Default() *Logger {
	return defaultLogger
}

// Debug 使用默认记录器记录调试级别日志
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Debugf 使用默认记录器记录格式化的调试级别日志
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Info 使用默认记录器记录信息级别日志
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof 使用默认记录器记录格式化的信息级别日志
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warn 使用默认记录器记录警告级别日志
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf 使用默认记录器记录格式化的警告级别日志
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Error 使用默认记录器记录错误级别日志
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf 使用默认记录器记录格式化的错误级别日志
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatal 使用默认记录器记录致命级别日志
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf 使用默认记录器记录格式化的致命级别日志
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Panic 使用默认记录器记录恐慌级别日志
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Panicf 使用默认记录器记录格式化的恐慌级别日志
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

// WithField 使用默认记录器添加单个字段
func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

// WithFields 使用默认记录器添加多个字段
func WithFields(fields map[string]interface{}) *Logger {
	return defaultLogger.WithFields(fields)
}
