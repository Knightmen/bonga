package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Resume struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	UserID    string    `json:"user_id" gorm:"not null"`
	RawText   string    `json:"raw_text" gorm:"type:text;not null"`
	Metadata  JSONB     `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type CreateResumeRequest struct {
	RawText  string `json:"raw_text" binding:"required"`
	Metadata JSONB  `json:"metadata"`
}

type ParseResumeRequest struct {
	FileName string `json:"fileName" binding:"required"`
}

type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, j)
} 