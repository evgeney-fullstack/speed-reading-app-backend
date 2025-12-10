package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCreateReadingTextIntegration is testing the endpoint of create a new reading text
func TestCreateReadingTextIntegration(t *testing.T) {
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
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "{\"id\":1}")
			},
		},
		{
			name: "Error -  invalid request body",
			payload: map[string]interface{}{
				"content":   "test test test test.",
				"questions": "question1",
			},
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "invalid_request_body", "Invalid request body format")
			},
		},
		{
			name: "Error - empty content",
			payload: map[string]interface{}{
				"content": "",
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
			contextTimeout: 15 * time.Second,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "invalid_request_body", "Invalid request body format")
			},
		},
		{
			name: "Error - request timeout",
			payload: map[string]interface{}{
				"content": strings.Repeat("test ", 1000000), // Большие данные
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
			contextTimeout: 100 * time.Millisecond,
			expectedStatus: http.StatusRequestTimeout,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "request_timeout", "Operation timed out")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request preparation
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/reading_text/", bytes.NewBuffer(body))
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
