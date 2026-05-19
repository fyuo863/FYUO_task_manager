package database

import (
	"FYUO_task_manager/internal/config"
	"FYUO_task_manager/pkg/log"
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(cfg *config.DatabaseConfig) error {
	// 使用 GORM 打开 MySQL 连接
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		// 配置 GORM 日志级别：生产环境可改为 Error 或 Silent
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	// 获取底层 sql.DB 对象，配置连接池参数（防止连接泄漏）
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       // 最大打开连接数
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 连接最大存活时间，到期后自动关闭重建

	// 使用带超时的 Context 测试数据库连通性
	ctx, cancel := context.WithTimeout(context.Background(), cfg.QueryTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("MySQL 连通性测试失败: %w", err)
	}
	DB = db
	log.Logger.Info("[DB] MySQL 连接池初始化完成")
	return nil
}

// Migrate 自动迁移传入的所有模型到 MySQL，即更新表结构
// models 为 GORM 模型结构体指针切片，如 &model.User{}, &model.Task{}
func Migrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化，请先调用 InitMySQL")
	}
	if err := DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("自动迁移表结构失败: %w", err)
	}
	log.Logger.Info("[DB] 表结构自动迁移完成")
	return nil
}

// func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
// 	return context.WithTimeout(parent, timeout)
// }
