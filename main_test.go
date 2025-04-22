package main_test

import (
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
	}{
		{
			name:          "positive case",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "email": "john@doe.com", "age": 18, "country": "USA"}`,
			wantCode:      http.StatusOK,
		},
		{
			name:          "positive case #2",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 5000, "email": "john@doe.com", "age": 130, "country": "Germany"}`,
			wantCode:      http.StatusOK,
		},
		{
			name:          "positive case #2",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 100000, "email": "john@doe.com", "age": 130, "country": "France"}`,
			wantCode:      http.StatusOK,
		},
		{
			name:          "non-specified ID",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"email": "john@doe.com", "age": 25, "country": "USA"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid ID",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": -1, "email": "john@doe.com", "age": 25, "country": "USA"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "non-specified email",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "age": 25, "country": "USA"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid email",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "email": "john.com", "age": 25, "country": "USA"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid age < 18",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "email": "john@doe.com", "age": 17, "country": "USA"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid age > 130",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "email": "john@doe.com", "age": 131, "country": "USA"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "invalid country",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "email": "john@doe.com", "age": 130, "country": "Unknown"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "non-specified age",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "email": "john@doe.com", "country": "USA"}`,
			wantCode:      http.StatusUnprocessableEntity,
		},
		{
			name:          "non-specified country",
			requestPath:   "/users",
			requestMethod: http.MethodPost,
			requestBody:   `{"id": 1, "email": "john@doe.com", "age": 18}`,
			wantCode:      http.StatusUnprocessableEntity,
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

			tr.NoError(tErr)
			tr.Equal(tc.wantCode, resp.StatusCode)
		})
	}
}
