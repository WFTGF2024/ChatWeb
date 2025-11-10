package models

import "time"

// ChatSession 映射 init.sql 中的 chat_sessions
type ChatSession struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	UserID    uint      `json:"user_id" gorm:"not null;column:user_id"`
	Title     string    `json:"title" gorm:"not null;column:title"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	// 若数据库中没有下列字段，请勿入库：仅在代码中使用时可解注并标记忽略
	// Model   string `json:"model" gorm:"-"`
	// Context string `json:"context" gorm:"-"`
}

func (ChatSession) TableName() string { return "chat_sessions" }

// ChatMessage 映射 init.sql 中的 chat_messages
type ChatMessage struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	UserID    uint      `json:"user_id" gorm:"not null;column:user_id"`
	SessionID uint      `json:"session_id" gorm:"not null;column:session_id"`
	Content   string    `json:"content" gorm:"type:text;column:content"`
	Role      string    `json:"role" gorm:"not null;column:role"` // 'user' | 'assistant'
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	// 仅本地使用的长文地址不入库
	// ContentURL string `json:"content_url" gorm:"-"`
}

func (ChatMessage) TableName() string { return "chat_messages" }
