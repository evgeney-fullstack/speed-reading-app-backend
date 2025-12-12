package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDeleteReadingTextIntegration is testing the endpoint of delete a reading text
func TestDeleteReadingTextIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()

	dbConfig, cleanup, err := setupTestContainer(ctx) //setupTestContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up test container: %v", err)
	}
	defer cleanup()

	router, err := setupTestServer(dbConfig)
	if err != nil {
		t.Fatalf("Failed to set up test server: %v", err)
	}

	tests := []struct {
		name           string
		payload        interface{}
		parameter      string
		contextTimeout time.Duration
		expectedStatus int
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful creation",
			payload: map[string]interface{}{
				"content": "test test test test.",
				"questions": []map[string]interface{}{
					{
						"question": "question1",
						"answer":   true,
					},
					{
						"question": "question2",
						"answer":   false,
					},
					{
						"question": "question3",
						"answer":   false,
					},
					{
						"question": "question4",
						"answer":   true,
					},
					{
						"question": "question5",
						"answer":   false,
					},
				},
			},
			parameter:      "1",
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "\"id\":1")
			},
		},
		{
			name:           "Successful delete a reading text by id",
			parameter:      "1",
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusNoContent,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "")
			},
		},
		{
			name:           "Error - invalid text_id param",
			parameter:      ":1",
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Invalid text ID format")
			},
		},
		{
			name:           "Error - invalid text_id value",
			parameter:      "0",
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Text ID must be positive")
			},
		},
		{
			name:           "Error - invalid text_id value",
			parameter:      "-1",
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Text ID must be positive")
			},
		},

		{
			name:           "Error - text not found",
			parameter:      "1",
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Text not found")
			},
		},
	}

	t.Run("Successful creation", func(t *testing.T) {
		// Request preparation
		body, _ := json.Marshal(tests[0].payload)
		req, _ := http.NewRequest("POST", "/reading_text/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Setting the context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), tests[0].contextTimeout)
		defer cancel()
		req = req.WithContext(ctx)

		// Request execution
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		// Checking the response status
		assert.Equal(t, tests[0].expectedStatus, recorder.Code)

		// Checking the response body
		if tests[0].checkResponse != nil {
			tests[0].checkResponse(t, recorder)
		}
	})

	for _, tt := range tests {
		if tt.name == "Successful creation" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			// Request preparation
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/reading_text/%s", tt.parameter), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Setting the context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.contextTimeout)
			defer cancel()
			req = req.WithContext(ctx)

			// Request execution
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			// Checking the response status
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// Checking the response body
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}
