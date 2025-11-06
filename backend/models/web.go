package models

import "time"

type ContentChunk struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`          // 所属用户
	WebPageID uint      `json:"web_page_id" gorm:"index"`      // 归属页面
	Content   string    `json:"content" gorm:"type:longtext"`  // 分块文本
	URL       string    `json:"url" gorm:"index"`              // 冗余便于查询
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// (user_id, url) 唯一，保证每个用户自己的库内不重复
type WebPage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`                       // 所属用户
	URL       string    `json:"url" gorm:"index:uniq_user_url,unique"`      // 配合 UserID 做唯一
	Title     string    `json:"title" gorm:"type:varchar(512)"`
	Content   string    `json:"content" gorm:"type:longtext"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
