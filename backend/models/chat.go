package models

import "time"

type ChatMessage struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null"`
    SessionID string    `json:"session_id" gorm:"not null"`
    Content   string    `json:"content" gorm:"not null"`
    Role      string    `json:"role" gorm:"not null"` // user/assistant
    CreatedAt time.Time `json:"created_at"`
}

type ChatSession struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null"`
    Title     string    `json:"title" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
