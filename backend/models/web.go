package models

import (
	"time"
)

type ContentChunk struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" gorm:"type:text"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WebPage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	URL       string    `json:"url" gorm:"unique"`
	Title     string    `json:"title"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
