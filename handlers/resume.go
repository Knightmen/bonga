package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"go-server/models"
	"go-server/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResumeHandler struct {
	db        *gorm.DB
	s3Service *services.S3Service
}

func NewResumeHandler(db *gorm.DB) *ResumeHandler {
	s3Service, err := services.NewS3Service()
	if err != nil {
		// Log the error but continue without S3 service
		// You might want to handle this differently based on your requirements
		return &ResumeHandler{db: db}
	}
	return &ResumeHandler{
		db:        db,
		s3Service: s3Service,
	}
}

// CreateResume godoc
// @Summary Create a new resume
// @Description Create a new resume by parsing a file through external service
// @Tags resume
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Param request body models.ParseResumeRequest true "Parse Resume Request"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/resume [post]
func (h *ResumeHandler) CreateResume(c *gin.Context) {
	var request models.ParseResumeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authorization token from environment
	authToken := os.Getenv("API_KEY")
	if authToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parse API token not configured"})
		return
	}

	// Make HTTP request to parse endpoint
	parseURL := fmt.Sprintf("http://localhost:8000/resume/parse?fileName=%s", url.QueryEscape(request.FileName))
	
	httpReq, err := http.NewRequest("GET", parseURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	
	// Set headers
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", authToken)

	// Make the request
	client := &http.Client{Timeout: 40 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to call parse API: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	// Parse the JSON response
	var parseResponse struct {
		TextContent  string `json:"text_content"`
		SessionId 	 string `json:"session_id"`
		Metadata 	 models.JSONB `json:"metadata"`
	}
	if err := json.Unmarshal(body, &parseResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Create resume record in database
	resume := models.Resume{
		UserID:  parseResponse.SessionId,
		RawText:  parseResponse.TextContent,
		Metadata: parseResponse.Metadata,
	}
	
	if err := h.db.Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume to database"})
		return
	}
	
	// Return the response
	c.JSON(resp.StatusCode, gin.H{
		"status_code": resp.StatusCode,
		"session_id": parseResponse.SessionId,
	})
}

// GetResume godoc
// @Summary Get a resume by ID
// @Description Get a resume by its ID
// @Tags resume
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Param id path string true "Resume ID"
// @Success 200 {object} models.Resume
// @Router /api/v1/resume/{id} [get]
func (h *ResumeHandler) GetResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.db.First(&resume, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	c.JSON(http.StatusOK, resume)
}

// UpdateResume godoc
// @Summary Update a resume
// @Description Update a resume by its ID
// @Tags resume
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Param id path string true "Resume ID"
// @Param resume body models.Resume true "Resume Data"
// @Success 200 {object} models.Resume
// @Router /api/v1/resume/{id} [put]
func (h *ResumeHandler) UpdateResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.db.First(&resume, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	if err := c.ShouldBindJSON(&resume); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resume.UpdatedAt = time.Now()

	if err := h.db.Save(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resume)
}

// DeleteResume godoc
// @Summary Delete a resume
// @Description Delete a resume by its ID
// @Tags resume
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Param id path string true "Resume ID"
// @Success 204 "No Content"
// @Router /api/v1/resume/{id} [delete]
func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	id := c.Param("id")
	var resume models.Resume

	if err := h.db.First(&resume, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	if err := h.db.Delete(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// LatestResume godoc
// @Summary Get latest resume
// @Description Get the most recently created resume
// @Tags resume
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Success 200 {object} models.Resume
// @Router /api/v1/resume/latest [get]
func (h *ResumeHandler) LatestResume(c *gin.Context) {
	var resume models.Resume

	if err := h.db.Order("id desc").First(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resume)
}

// GetSignedURL godoc
// @Summary Get a presigned URL for uploading a resume
// @Description Get a presigned URL for uploading a resume to S3
// @Tags resume
// @Accept json
// @Produce json
// @Param Authorization header string true "API Key"
// @Param filename query string true "Filename for the resume"
// @Success 200 {object} map[string]string
// @Router /api/v1/resume/getSignedUrl [get]
func (h *ResumeHandler) GetSignedURL(c *gin.Context) {
	if h.s3Service == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "S3 service is not available"})
		return
	}

	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	key := "resumes/" + filename

	// Get presigned URL
	url, err := h.s3Service.GetPresignedURL(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
		"key": key,
	})
} 