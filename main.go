package main

import (
	"FYUO_task_manager/pkg/log"
	"fmt"

	"FYUO_task_manager/internal/config"
	"FYUO_task_manager/internal/router"
)

func main() {
	log.LogSetting()
	cfg, err := config.Load("configs/config.yaml") // cfg 为读取到的配置
	if err != nil {
		log.Logger.Fatalf("[启动失败] 加载配置: %v", err)
	}
	log.Logger.Printf("[Config] 端口号: %s, 运行模式: %s", cfg.Server.Port, cfg.Server.Mode)
	router.Setup(cfg.Server.Mode)

	defer func() { fmt.Println("安全退出占位") }()
}
