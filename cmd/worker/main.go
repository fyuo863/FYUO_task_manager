package main

import (
	"FYUO_task_manager/internal/config"
	"FYUO_task_manager/internal/database"
	"FYUO_task_manager/internal/model"
	"FYUO_task_manager/internal/worker"
	"FYUO_task_manager/pkg/log"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.LogSetting()
	cfg, err := config.Load("configs/config.yaml") // cfg 为读取到的配置
	if err != nil {
		log.Logger.Fatalf("[启动失败] 加载配置: %v", err)
	}
	fmt.Println(cfg)

	//初始化database
	if err := database.InitMySQL(&cfg.Database); err != nil {
		log.Logger.Fatalf("[Worker启动失败] 初始化 MySQL: %v", err)
	}
	log.Logger.Info("[Worker]", "MySQL", "连接成功")
	// 自动建表：根据 GORM 模型的 struct tag 创建/更新表结构
	if err := database.Migrate(&model.User{}, &model.Task{}); err != nil {
		log.Logger.Fatalf("[启动失败] 数据库迁移: %v", err)
	}

	//初始化Redis

	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Logger.Fatalf("[启动失败] 初始化Redis: %v", err)
	}
	defer database.CloseRedis() //关闭Redis连接

	log.Logger.Info("[Config]", "端口号", cfg.Server.Port, "运行模式", cfg.Server.Mode)

	pool := worker.NewPool(5, 100) //新建worker pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT) //监测是否触发Ctrl+C

	go func() {
		log.Logger.Info("[Worker] 主循环已启动，等待任务队列...")

		pool.ListenQueue(ctx, database.RDB)

	}()

	sig := <-quit //阻塞Http协程直至触发退出
	log.Logger.Info("[Server]", "接收到退出信号", sig)
	defer func() {
		log.Logger.Error("安全退出占位")
		pool.Stop()
	}()
}
