package models

import (
	"time"
)

// ReadingText presents the text to check the reading speed
type ReadingText struct {
	ID        int64     `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	WordCount int       `json:"word_count" db:"word_count" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReadingQuestion presents a question to the text
type ReadingQuestion struct {
	ID        int64     `json:"id" db:"id"`
	TextID    int64     `json:"text_id" db:"text_id"`
	Question  string    `json:"question" db:"question" binding:"required"`
	Answer    bool      `json:"answer" db:"answer" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReadingSessionResult represents the result of a reading session
type ReadingSessionResult struct {
	WordsPerMinute float64 `json:"words_per_minute"`
	Comprehension  float64 `json:"comprehension_percent"` // 0–100
}
