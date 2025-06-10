package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// Prepare request body
	reqBody := map[string]string{
		"session_id": request.SessionID,
		"message":   request.Question,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	// Make POST request to chat endpoint
	resp, err := http.Post("http://localhost:8000/session/chat", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to call chat API: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Parse response
	var chatResponse struct {
		Answer    string `json:"answer"`
		SessionID string `json:"session_id"`
		Status    string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Return chat response
	response := SessionChatResponse{
		SessionID: chatResponse.SessionID,
		Answer:    chatResponse.Answer,
	}

	c.JSON(http.StatusOK, response)
}