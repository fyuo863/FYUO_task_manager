package database

import (
	"FYUO_task_manager/internal/config"
	"FYUO_task_manager/pkg/log"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(), // Redis 服务器地址
		Password: cfg.Password,    // Redis 密码（无密码则为空字符串）
		DB:       cfg.DB,          // Redis 数据库编号
		PoolSize: cfg.PoolSize,    // 连接池大小，高并发场景防止连接不够用
	})
	// 2. 设置超时上下文，验证连接是否成功
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	con, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}
	RDB = rdb
	log.Logger.Info("[Config]", "Redis初始化成功", con)
	return nil
}

func CloseRedis() {
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Logger.Info("[Redis]", "关闭连接时发生错误", err)
			return
		}
		log.Logger.Info("[Redis] 连接已正常关闭")
	}
}
