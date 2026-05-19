package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"` // 服务监听地址
	Port int    `mapstructure:"port"` // 服务监听端口
	Mode string `mapstructure:"mode"` // Gin 运行模式: debug / release / test
}

type DatabaseConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`            // Redis 数据库编号
	PoolSize     int           `mapstructure:"pool_size"`     // 连接池大小
	QueryTimeout time.Duration `mapstructure:"query_timeout"` // 单次操作超时时间
}

func Load(configPath string) (*Config, error) {
	//使用viper加载配置文件
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil { //将配置文件内容读取到viper内部配置表中
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil { //将配置内容存入cfg结构体中
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}

// Addr 返回 Redis 连接地址
// func (s *ServerConfig) Serve() *ServerConfig        { return s }
// func (r *RedisConfig) Redis() *RedisConfig          { return r }
// func (r *DatabaseConfig) Database() *DatabaseConfig { return r }
//
//	func (r *RedisConfig) Addr() string {
//		return fmt.Sprintf("%s:%d", r.Host, r.Port)
//	}
func (s *ServerConfig) ServeAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
func (r *RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
func (d *DatabaseConfig) DatabaseAddr() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}
