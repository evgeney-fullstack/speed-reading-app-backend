package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ReadingText presents the text to check the reading speed
type ReadingText struct {
	ID        int64        `json:"id" db:"id"`
	Content   string       `json:"content" db:"content" binding:"required,min=1"`
	WordCount int          `json:"word_count" db:"word_count"`
	Questions QuestionList `json:"questions" db:"questions" binding:"required,min=1"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

// Question presents a question to the text
type Question struct {
	ID       int64  `json:"id,omitempty"`
	Question string `json:"question" binding:"required"`
	Answer   bool   `json:"answer" binding:"required"`
}

// QuestionList is a custom type for working with JSON in PostgreSQL
type QuestionList []Question

// Value - serialization for saving to the database
func (q QuestionList) Value() (driver.Value, error) {
	if q == nil {
		return "[]", nil
	}
	return json.Marshal(q)
}

// Scan deserialization when reading from a database
func (q *QuestionList) Scan(value interface{}) error {
	if value == nil {
		*q = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, q)
}
