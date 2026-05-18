package main

import (
	"FYUO_task_manager/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"FYUO_task_manager/internal/config"
	"FYUO_task_manager/internal/router"
)

func main() {
	log.LogSetting()
	cfg, err := config.Load("configs/config.yaml") // cfg 为读取到的配置
	if err != nil {
		log.Logger.Fatalf("[启动失败] 加载配置: %v", err)
	}
	log.Logger.Info("[Config]", "端口号", cfg.Server.Port, "运行模式", cfg.Server.Mode)
	r := router.Setup(cfg.Server.Mode)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r, // 将 Gin 引擎作为 HTTP Handler 传入
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT) //监测是否触发Ctrl+C

	go func() { //在协程中启动Http服务
		log.Logger.Info("[Server]", "HTTP 服务启动于", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Fatal("[Server]", "服务启动失败, err", err)
		}
	}()

	sig := <-quit //阻塞Http协程直至触发退出
	log.Logger.Info("[Server]", "接收到退出信号", sig)
	defer func() { log.Logger.Error("安全退出占位") }()
}
