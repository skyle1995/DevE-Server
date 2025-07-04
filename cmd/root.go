package cmd

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/skyle1995/DevE-Server/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd 表示没有调用子命令时的基础命令
var rootCmd = &cobra.Command{
	Use:   "DevE",
	Short: "DevE网络验证系统",
	Long: `DevE网络验证系统，以Golang为主语言开发，支持多用户，多应用。
该系统主要包含用户后台核心模块，支持多应用管理，每个应用拥有独立的卡密和设备管理。`,
	// 如果有任何初始化操作需要在每个子命令之前运行，可以放在这里
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 初始化日志
		if logFile != "" {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
			if err != nil {
				log.Fatal("无法创建日志文件:", err)
			}
			mw := io.MultiWriter(os.Stdout, file)
			log.SetOutput(mw)
		} else {
			mw := io.MultiWriter(os.Stdout)
			log.SetOutput(mw)
		}

		if logType == "json" {
			log.SetFormatter(&log.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
			})
		} else {
			log.SetFormatter(&log.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
			})
		}

		// 调用config包中的初始化函数
		config.Init(configFile)

		if port != 0 {
			viper.Set("server.port", port)
		}

		if debug {
			viper.Set("server.level", "DEBUG")
		}

		if viper.IsSet("server.level") {
			logLevel := viper.GetString("server.level")
			switch logLevel {
			case "DEBUG":
				log.SetLevel(log.DebugLevel)
			case "INFO":
				log.SetLevel(log.InfoLevel)
			case "WARN":
				log.SetLevel(log.WarnLevel)
			case "ERROR":
				log.SetLevel(log.ErrorLevel)
			case "FATAL":
				log.SetLevel(log.FatalLevel)
			default:
				log.Info("未知的日志等级，使用 INFO 等级")
				log.SetLevel(log.InfoLevel)
			}
			log.WithFields(log.Fields{"level": logLevel}).Info("使用日志等级")
		} else {
			log.SetLevel(log.InfoLevel)
			log.Info("未设置日志等级，使用默认等级")
		}
	},
}

// Execute 将所有子命令添加到root命令并适当设置标志。
// 这由main.main()调用。只需要对rootCmd调用一次。
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Errorf("执行命令出错: %v", err)
		os.Exit(1)
	}
}

// init 初始化命令行参数和子命令
func init() {
	// 日志文件
	rootCmd.PersistentFlags().StringVar(&logFile, "log", "", "日志文件路径 默认为空，不输出日志文件")
	// 日志文件
	rootCmd.PersistentFlags().StringVar(&logType, "type", "text", "日志输出格式 (默认为text 可选：text/json)")
	// 配置文件路径
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "配置文件路径 (默认为 config.yaml)")
	// 服务器端口
	rootCmd.PersistentFlags().IntVar(&port, "port", 0, "服务器端口")
	// 调试模式
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "启用调试模式")

	// 添加子命令
	rootCmd.AddCommand(startCmd)
}
