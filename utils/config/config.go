package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config 配置管理器结构体
type Config struct {
	viper *viper.Viper
	path  string
}

// New 创建新的配置管理器
func New() *Config {
	return &Config{
		viper: viper.New(),
	}
}

// LoadFile 从文件加载配置
func (c *Config) LoadFile(path string) error {
	ext := filepath.Ext(path)
	c.viper.SetConfigFile(path)
	c.path = path

	switch strings.ToLower(ext) {
	case ".json":
		c.viper.SetConfigType("json")
	case ".yaml", ".yml":
		c.viper.SetConfigType("yaml")
	case ".toml":
		c.viper.SetConfigType("toml")
	case ".ini":
		c.viper.SetConfigType("ini")
	case ".env":
		c.viper.SetConfigType("env")
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	err := c.viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

// LoadJSON 从JSON字符串加载配置
func (c *Config) LoadJSON(jsonStr string) error {
	c.viper.SetConfigType("json")
	return c.viper.ReadConfig(strings.NewReader(jsonStr))
}

// LoadYAML 从YAML字符串加载配置
func (c *Config) LoadYAML(yamlStr string) error {
	c.viper.SetConfigType("yaml")
	return c.viper.ReadConfig(strings.NewReader(yamlStr))
}

// Get 获取配置项
func (c *Config) Get(key string) interface{} {
	return c.viper.Get(key)
}

// GetString 获取字符串配置项
func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetInt 获取整数配置项
func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetInt64 获取64位整数配置项
func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

// GetFloat64 获取浮点数配置项
func (c *Config) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// GetBool 获取布尔配置项
func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetStringSlice 获取字符串切片配置项
func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

// GetStringMap 获取字符串映射配置项
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}

// GetStringMapString 获取字符串映射字符串配置项
func (c *Config) GetStringMapString(key string) map[string]string {
	return c.viper.GetStringMapString(key)
}

// Set 设置配置项
func (c *Config) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

// Has 检查配置项是否存在
func (c *Config) Has(key string) bool {
	return c.viper.IsSet(key)
}

// AllSettings 获取所有配置项
func (c *Config) AllSettings() map[string]interface{} {
	return c.viper.AllSettings()
}

// AllKeys 获取所有配置键
func (c *Config) AllKeys() []string {
	return c.viper.AllKeys()
}

// Save 保存配置到文件
func (c *Config) Save() error {
	if c.path == "" {
		return fmt.Errorf("config file path not set")
	}

	var (
		data []byte
		err  error
	)

	ext := filepath.Ext(c.path)
	switch strings.ToLower(ext) {
	case ".json":
		data, err = json.MarshalIndent(c.viper.AllSettings(), "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(c.viper.AllSettings())
	default:
		return fmt.Errorf("unsupported config file format for saving: %s", ext)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0644)
}

// SaveAs 保存配置到指定文件
func (c *Config) SaveAs(path string) error {
	oldPath := c.path
	c.path = path
	err := c.Save()
	if err != nil {
		c.path = oldPath
		return err
	}
	return nil
}

// SetDefault 设置默认配置项
func (c *Config) SetDefault(key string, value interface{}) {
	c.viper.SetDefault(key, value)
}

// SetEnvPrefix 设置环境变量前缀
func (c *Config) SetEnvPrefix(prefix string) {
	c.viper.SetEnvPrefix(prefix)
}

// AutomaticEnv 启用自动环境变量绑定
func (c *Config) AutomaticEnv() {
	c.viper.AutomaticEnv()
}

// BindEnv 绑定环境变量
func (c *Config) BindEnv(key string, envVar ...string) error {
	args := append([]string{key}, envVar...)
	return c.viper.BindEnv(args...)
}

// Unmarshal 将配置解析到结构体
func (c *Config) Unmarshal(rawVal interface{}) error {
	return c.viper.Unmarshal(rawVal)
}

// UnmarshalKey 将配置键解析到结构体
func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	return c.viper.UnmarshalKey(key, rawVal)
}

// WatchConfig 监视配置文件变化
func (c *Config) WatchConfig() {
	c.viper.WatchConfig()
}

// OnConfigChange 设置配置变化回调函数
func (c *Config) OnConfigChange(run func()) {
	c.viper.OnConfigChange(func(_ fsnotify.Event) {
		run()
	})
}

// LoadDefault 加载默认配置文件
func LoadDefault() (*Config, error) {
	config := New()

	// 尝试按优先级加载配置文件
	configPaths := []string{
		"./config.yaml",
		"./config.yml",
		"./config.json",
		"./config/config.yaml",
		"./config/config.yml",
		"./config/config.json",
		"../config.yaml",
		"../config.yml",
		"../config.json",
		"../config/config.yaml",
		"../config/config.yml",
		"../config/config.json",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			err := config.LoadFile(path)
			if err == nil {
				return config, nil
			}
		}
	}

	return nil, fmt.Errorf("no config file found")
}
