package services

import (
	"backend/models"
    "backend/database"
)

type ChatService struct{}

// NewChatService 创建新的聊天服务实例
func NewChatService() *ChatService {
    return &ChatService{}
}

func (s *ChatService) CreateSession(userID uint, title string) (*models.ChatSession, error) {
    session := &models.ChatSession{
        UserID: userID,
        Title:  title,
    }
    if err := database.DB.Create(session).Error; err != nil {
        return nil, err
    }
    return session, nil
}

func (s *ChatService) SaveMessage(userID uint, sessionID string, content, role string) error {
    message := &models.ChatMessage{
        UserID:    userID,
        SessionID: sessionID,
        Content:   content,
        Role:      role,
    }
    return database.DB.Create(message).Error
}

func (s *ChatService) GetChatHistory(userID uint) ([]models.ChatSession, error) {
    var sessions []models.ChatSession
    err := database.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&sessions).Error
    return sessions, err
}

func (s *ChatService) GetSessionMessages(userID uint, sessionID string) ([]models.ChatMessage, error) {
    var messages []models.ChatMessage
    err := database.DB.Where("user_id = ? AND session_id = ?", userID, sessionID).
        Order("created_at ASC").Find(&messages).Error
    return messages, err
}
