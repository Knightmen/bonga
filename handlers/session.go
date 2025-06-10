package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SessionHandler struct {
	db *gorm.DB
}

func NewSessionHandler(db *gorm.DB) *SessionHandler {
	return &SessionHandler{
		db: db,
	}
}

// SessionChatRequest represents the request body for session chat
type SessionChatRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Question  string `json:"question" binding:"required"`
}

// SessionChatResponse represents the response for session chat
type SessionChatResponse struct {
	SessionID string `json:"sessionId"`
	Answer    string `json:"answer"`
}

// InitSession godoc
// @Summary Initialize a new session
// @Description Initialize a new chat session
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/session/init [get]
func (h *SessionHandler) InitSession(c *gin.Context) {
	// Static/dummy response for now
	c.JSON(http.StatusOK, gin.H{
		"sessionId": "session-12345",
		"status":    "initialized",
		"message":   "Session initialized successfully",
	})
}

// ChatSession godoc
// @Summary Chat with a session
// @Description Send a question to a chat session and get an answer
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Param request body SessionChatRequest true "Chat Request"
// @Success 200 {object} SessionChatResponse
// @Router /api/v1/session/chat [post]
func (h *SessionHandler) ChatSession(c *gin.Context) {
	var request SessionChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Static/dummy response for now
	response := SessionChatResponse{
		SessionID: request.SessionID,
		Answer:    "This is a dummy response to your question: " + request.Question,
	}

	c.JSON(http.StatusOK, response)
} 