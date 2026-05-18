package main

import (
	"log"

	"FYUO_task_manager/internal/config"
)

func main() {
	cfg, err := config.Load("configs/config.yaml") // cfg 为读取到的配置
	if err != nil {
		log.Fatalf("[启动失败] 加载配置: %v", err)
	}
	log.Printf("[Config] 配置加载成功 — 服务端口: %s, 运行模式: %s", cfg.Server.Port, cfg.Server.Mode)

}
