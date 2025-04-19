package main_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	testCases := []struct {
		name          string
		requestPath   string
		requestBody   string
		requestMethod string
		wantCode      int
		wantBody      string
	}{
		{
			name:          "Task not found",
			requestPath:   "/tasks/1",
			requestMethod: http.MethodGet,
			wantCode:      http.StatusNotFound,
			wantBody:      "Not Found",
		},
		{
			name:          "Create a new task",
			requestPath:   "/tasks",
			requestMethod: http.MethodPost,
			requestBody:   `{"description": "Learn Go", "deadline": 1620000000}`,
			wantCode:      http.StatusOK,
			wantBody:      `{"id":1}`,
		},
		{
			name:          "Get the newly created task",
			requestPath:   "/tasks/1",
			requestMethod: http.MethodGet,
			wantCode:      http.StatusOK,
			wantBody:      `{"id":1,"description":"Learn Go","deadline":1620000000}`,
		},
		{
			name:          "Update the task #1",
			requestPath:   "/tasks/1",
			requestMethod: http.MethodPatch,
			requestBody:   `{"description": "Learn Go and Fiber", "deadline": 1630000000}`,
			wantCode:      http.StatusOK,
			wantBody:      `OK`,
		},
		{
			name:          "Get the updated task #1",
			requestPath:   "/tasks/1",
			requestMethod: http.MethodGet,
			wantCode:      http.StatusOK,
			wantBody:      `{"id":1,"description":"Learn Go and Fiber","deadline":1630000000}`,
		},
		{
			name:          "Delete the task #1",
			requestPath:   "/tasks/1",
			requestMethod: http.MethodDelete,
			wantCode:      http.StatusOK,
			wantBody:      `OK`,
		},
		{
			name:          "Task #1 not exists",
			requestPath:   "/tasks/1",
			requestMethod: http.MethodGet,
			wantCode:      http.StatusNotFound,
			wantBody:      "Not Found",
		},
		{
			name:          "Create a new task #2",
			requestPath:   "/tasks",
			requestMethod: http.MethodPost,
			requestBody:   `{"description": "Learn Go #2", "deadline": 1620000000}`,
			wantCode:      http.StatusOK,
			wantBody:      `{"id":2}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, tErr := http.NewRequest(
				tc.requestMethod,
				"http://localhost:8080"+tc.requestPath,
				strings.NewReader(tc.requestBody),
			)
			tr := require.New(t)
			tr.NoError(tErr)

			req.Header.Set("Content-Type", "application/json")

			httpClient := http.Client{}
			resp, tErr := httpClient.Do(req)
			tr.NoError(tErr)
			defer resp.Body.Close()

			bodyBytes, tErr := io.ReadAll(resp.Body)
			tr.NoError(tErr)

			tr.Equal(tc.wantCode, resp.StatusCode)
			tr.Equal(tc.wantBody, string(bodyBytes))
		})
	}
}
