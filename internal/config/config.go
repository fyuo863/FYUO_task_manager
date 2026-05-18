package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"` // 服务监听端口
	Mode string `mapstructure:"mode"` // Gin 运行模式: debug / release / test
}

type DatabaseConfig struct {
}

type RedisConfig struct {
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
