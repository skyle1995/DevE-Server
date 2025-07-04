package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Server 表示HTTP服务器
type Server struct {
	httpServer *http.Server
}

// NewServer 创建一个新的服务器实例
func NewServer() *Server {
	// 获取配置的端口
	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8025 // 默认端口
	}

	// 设置路由
	router := SetupRouter()

	// 创建HTTP服务器
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &Server{
		httpServer: httpServer,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop 优雅地关闭服务器
func (s *Server) Stop() error {
	log.Info("正在关闭服务器...")

	// 创建一个5秒的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Errorf("服务器关闭出错: %v", err)
		return err
	}

	log.Info("服务器已成功关闭")
	return nil
}
