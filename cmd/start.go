package cmd

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/skyle1995/DevE-Server/database"
	"github.com/skyle1995/DevE-Server/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd 表示启动服务器的命令
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动DevE网络验证系统服务器",
	Long: `启动DevE网络验证系统服务器，并开始监听HTTP请求。
该命令将初始化数据库连接，加载配置，并启动Web服务器。
服务器将一直运行，直到收到终止信号。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化数据库
		dbType := viper.GetString("database.type")
		log.Infof("正在初始化%s数据库连接...", dbType)

		// 初始化数据库连接
		database.Init()

		// 执行数据库迁移
		migration := database.NewMigration()
		if err := migration.AutoMigrate(); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}

		// 初始化默认数据
		if err := migration.InitDefaultData(); err != nil {
			log.Fatalf("初始化默认数据失败: %v", err)
		}

		log.Info("数据库初始化完成")

		// --------------------------------------------------------- //

		// 启动服务器
		log.Infof("正在启动服务器，监听端口: %d", viper.GetInt("server.port"))

		// 创建并启动服务器
		srv := server.NewServer()
		go func() {
			if err := srv.Start(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("服务器启动失败: %v", err)
			}
		}()

		log.Infof("服务器已启动，访问地址: http://localhost:%d", viper.GetInt("server.port"))

		// --------------------------------------------------------- //

		// 等待终止信号
		// 创建一个通道来接收信号
		sigChan := make(chan os.Signal, 1)
		// 监听终止信号
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		// 等待信号
		sig := <-sigChan
		log.Infof("收到信号: %v，正在优雅关闭服务器...", sig)

		// 关闭服务器
		if err := srv.Stop(); err != nil {
			log.Errorf("服务器关闭出错: %v", err)
		}

		// 关闭数据库连接
		database.Close()
		log.Info("数据库连接已关闭")

		log.Info("服务器已关闭")
	},
}

func init() {
	// 添加启动命令特定的标志
	rootCmd.AddCommand(startCmd)
}
