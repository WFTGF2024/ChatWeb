// backend/services/chat_service.go
package services

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"backend/database"
	"backend/models"
)

type ChatService struct{}

// NewChatService 创建新的聊天服务实例
func NewChatService() *ChatService {
	return &ChatService{}
}

// CreateSession 创建会话（title 为空则给默认值）
func (s *ChatService) CreateSession(userID uint, title string) (*models.ChatSession, error) {
	if strings.TrimSpace(title) == "" {
		title = "新会话"
	}
	session := &models.ChatSession{
		UserID: userID,
		Title:  strings.TrimSpace(title),
	}
	if err := database.DB.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

// SaveMessage 追加一条消息：把字符串 sessionID 转成 uint，入库后刷新会话的 updated_at
// 说明：不再使用 ContentURL（模型/表中无此字段），如需大文本外链存储，请先在表与模型中新增列再改这里。
func (s *ChatService) SaveMessage(userID uint, sessionID, content, role string) error {
	if strings.TrimSpace(content) == "" {
		return errors.New("content required")
	}
	if role == "" {
		role = "user"
	}

	sid, err := parseSessionID(sessionID)
	if err != nil {
		return err
	}

	message := models.ChatMessage{
		UserID:    userID,
		SessionID: sid,                 // ✅ 转为 uint
		Content:   strings.TrimSpace(content), // ✅ 仅存内容本身
		Role:      role,
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(&message).Error; err != nil {
		return err
	}

	// 刷新会话的更新时间（若表中存在 updated_at）
	_ = database.DB.
		Model(&models.ChatSession{}).
		Where("id = ? AND user_id = ?", sid, userID).
		Update("updated_at", message.CreatedAt).Error

	return nil
}

// GetFullMessage 返回消息内容本身（当前不做外链/压缩恢复）
func (s *ChatService) GetFullMessage(message *models.ChatMessage) (string, error) {
	if message == nil {
		return "", errors.New("nil message")
	}
	return message.Content, nil
}

// GetChatHistory 获取用户的会话列表（按更新时间倒序，其次按 id 倒序）
func (s *ChatService) GetChatHistory(userID uint) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	err := database.DB.
		Where("user_id = ?", userID).
		Order("updated_at DESC, id DESC").
		Find(&sessions).Error
	return sessions, err
}

// GetSessionMessages 获取某会话的消息（按时间升序，其次按 id 升序）
func (s *ChatService) GetSessionMessages(userID uint, sessionID string) ([]models.ChatMessage, error) {
	sid, err := parseSessionID(sessionID)
	if err != nil {
		return nil, err
	}
	var messages []models.ChatMessage
	err = database.DB.
		Where("user_id = ? AND session_id = ?", userID, sid).
		Order("created_at ASC, id ASC").
		Find(&messages).Error
	return messages, err
}

//
// ===== 供 handlers 直接调用的包级函数（与历史代码兼容） =====
//

// AddMessage 追加一条消息（handlers 依赖的导出函数）
func AddMessage(userID uint, sessionIDStr, role, content string) (*models.ChatMessage, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("content required")
	}
	if role == "" {
		role = "user"
	}

	sid, err := parseSessionID(sessionIDStr)
	if err != nil {
		return nil, err
	}

	msg := &models.ChatMessage{
		UserID:    userID,
		SessionID: sid,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}
	if err := database.DB.Create(msg).Error; err != nil {
		return nil, err
	}

	// 刷新会话 updated_at
	_ = database.DB.
		Model(&models.ChatSession{}).
		Where("id = ? AND user_id = ?", sid, userID).
		Update("updated_at", msg.CreatedAt).Error

	return msg, nil
}

// ListMessages 拉取会话消息（按时间升序；limit<=0 表示不限）
func ListMessages(userID uint, sessionIDStr string, limit int) ([]models.ChatMessage, error) {
	sid, err := parseSessionID(sessionIDStr)
	if err != nil {
		return nil, err
	}
	var out []models.ChatMessage
	q := database.DB.
		Where("user_id = ? AND session_id = ?", userID, sid).
		Order("created_at ASC, id ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

//
// --- helpers ---
//

func parseSessionID(s string) (uint, error) {
	if s == "" {
		return 0, errors.New("empty session_id")
	}
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}
