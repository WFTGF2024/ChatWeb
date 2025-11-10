package handlers

import (
	"net/http"
    "backend/services"
	"github.com/gin-gonic/gin"
)

func HandleCreateSession(c *gin.Context) {
    userID := c.GetUint("user_id")
    var req struct {
        Title string `json:"title"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    service := services.NewChatService()
    session, err := service.CreateSession(userID, req.Title)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, session)
}

func HandleSaveMessage(c *gin.Context) {
    userID := c.GetUint("user_id")
    var req struct {
        SessionID string `json:"session_id"`
        Content   string `json:"content"`
        Role      string `json:"role"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    service := services.NewChatService()
    if err := service.SaveMessage(userID, req.SessionID, req.Content, req.Role); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": true})
}

func HandleGetSessionMessages(c *gin.Context) {
    userID := c.GetUint("user_id")
    sessionID := c.Param("session_id")
    
    service := services.NewChatService()
    messages, err := service.GetSessionMessages(userID, sessionID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, messages)
}
