package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/library-system/backend/internal/config"
	"github.com/library-system/backend/internal/routes"
	"github.com/library-system/backend/pkg/database"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接 MySQL
	mysqlDB, err := database.InitMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect MySQL: %v", err)
	}
	defer database.CloseMySQL()

	// 连接 MongoDB
	mongoDB, err := database.InitMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect MongoDB: %v", err)
	}
	defer database.CloseMongoDB()

	// 设置路由
	router := routes.SetupRouter(mysqlDB, mongoDB, cfg)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("🚀 Server starting on port %d...", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("👋 Server exited")
}
