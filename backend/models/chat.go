package models

import (
	"time"
)

type ChatHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileBridge struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FileID    string    `json:"file_id"`
	UserID    uint      `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
