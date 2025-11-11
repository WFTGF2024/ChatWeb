// backend/handlers/chat_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"backend/database"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

//
// 会话 CRUD
//

// POST /api/chats  创建会话
func HandleCreateSession(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}
	row := &models.ChatSession{UserID: userID, Title: strings.TrimSpace(req.Title)}
	if err := database.DB.Create(row).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, row)
}

// GET /api/chats  列出我的会话
func HandleListSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	var items []models.ChatSession
	if err := database.DB.
		Where("user_id = ?", userID).
		Order("updated_at DESC, id DESC").
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

//
// 消息 CRUD
//

// GET /api/chats/:session_id/messages  获取某会话消息
func HandleGetSessionMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	sess := c.Param("session_id")
	id64, err := strconv.ParseUint(sess, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}
	var messages []models.ChatMessage
	if err := database.DB.
		Where("user_id = ? AND session_id = ?", userID, uint(id64)).
		Order("created_at ASC, id ASC").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// POST /api/chats/:session_id/messages  追加一条消息
func HandleAddMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	sess := c.Param("session_id")
	var req struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content required"})
		return
	}
	if req.Role == "" {
		req.Role = "user"
	}
	msg, err := services.AddMessage(userID, sess, req.Role, strings.TrimSpace(req.Content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

// ✅ 兼容别名：有些 main.go 写的是 HandleSaveMessage，这里做一个别名避免你改路由
// 支持两种入参：路径参数 :session_id 或 JSON 里的 session_id
// POST /api/chats/:session_id/messages   或   POST /api/chats/messages {session_id, role, content}
func HandleSaveMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 先尝试从 path 取
	sess := c.Param("session_id")

	// 再尝试 JSON
	var req struct {
		SessionID string `json:"session_id"`
		Role      string `json:"role"`
		Content   string `json:"content"`
	}
	_ = c.ShouldBindJSON(&req)

	if sess == "" {
		sess = req.SessionID
	}
	if strings.TrimSpace(sess) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}
	role := req.Role
	if role == "" {
		role = "user"
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content required"})
		return
	}

	msg, err := services.AddMessage(userID, sess, role, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

//
// 上下文 + LLM
//

// 内部：构建上下文（最近 N 条）
func buildContextMessages(userID, sessionID uint, limit int) ([]services.LLMMessage, error) {
	var msgs []models.ChatMessage
	q := database.DB.
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("created_at ASC, id ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&msgs).Error; err != nil {
		return nil, err
	}
	out := make([]services.LLMMessage, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, services.LLMMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return out, nil
}

// POST /api/chats/:session_id/complete  一次性补全（可携带 content，先写 user 再补全）
func HandleLLMCompleteOnce(c *gin.Context) {
	userID := c.GetUint("user_id")
	sess := c.Param("session_id")
	id64, err := strconv.ParseUint(sess, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}
	sessionID := uint(id64)

	var req struct {
		Content string `json:"content"`
		Model   string `json:"model"`
	}
	_ = c.ShouldBindJSON(&req)

	if strings.TrimSpace(req.Content) != "" {
		if _, err := services.AddMessage(userID, sess, "user", strings.TrimSpace(req.Content)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	messages, err := buildContextMessages(userID, sessionID, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	llm := services.NewLLMClient()
	text, err := llm.ChatOnce(c.Request.Context(), messages, req.Model)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// 保存助手消息（失败不影响返回）
	if _, _err := services.AddMessage(userID, sess, "assistant", text); _err != nil {
		// 可记录日志
	}

	c.JSON(http.StatusOK, gin.H{"content": text})
}

// POST /api/chats/:session_id/stream  流式补全
func HandleLLMStream(c *gin.Context) {
	userID := c.GetUint("user_id")
	sess := c.Param("session_id")
	id64, err := strconv.ParseUint(sess, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}
	sessionID := uint(id64)

	var req struct {
		Content string `json:"content"`
		Model   string `json:"model"`
	}
	_ = c.ShouldBindJSON(&req)

	if strings.TrimSpace(req.Content) != "" {
		if _, err := services.AddMessage(userID, sess, "user", strings.TrimSpace(req.Content)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	messages, err := buildContextMessages(userID, sessionID, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	llm := services.NewLLMClient()
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Status(http.StatusOK)
	_ = llm.ChatStream(
		c.Request.Context(),
		c.Writer,
		messages,
		req.Model,
		func(full string) {
			// 流结束后保存助手消息（失败也不影响客户端）
			_, _ = services.AddMessage(userID, sess, "assistant", full)
		},
	)
}

// DELETE /api/chats/:session_id  删除会话
func HandleDeleteSession(c *gin.Context) {
    userID := c.GetUint("user_id")
    sessionID := c.Param("session_id")
    
    if _, err := strconv.ParseUint(sessionID, 10, 64); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
        return
    }

    // 使用指针调用方法
    chatService := &services.ChatService{}
    if err := chatService.DeleteSession(userID, sessionID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "message": "会话已删除"})
}
