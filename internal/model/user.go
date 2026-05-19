package model

import "time"

// User 用户表 — 对应 MySQL users 表
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`                    // 用户主键
	Username  string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"` // 登录用户名（唯一）
	Password  string    `gorm:"type:varchar(256);not null" json:"-"`                   // 密码哈希（json:"-" 防止序列化泄露）
	Email     string    `gorm:"type:varchar(128)" json:"email"`                        // 邮箱
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                      // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`                      // 更新时间
}
