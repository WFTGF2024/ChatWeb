package services

import (
	"backend/models"
)

type ChatService struct{}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) GetChatHistory(userID int) ([]models.ChatHistory, error) {
	// 实现获取聊天历史逻辑
	return nil, nil
}

func (s *ChatService) BridgeFileServer(fileID string) error {
	// 实现文件服务器桥接逻辑
	return nil
}
