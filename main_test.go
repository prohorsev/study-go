package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	testCases := []struct {
		name          string
		postID        string
		requestMethod string
		wantCode      int
		want          string
	}{
		{
			name:          "post not found",
			postID:        "12345",
			requestMethod: fiber.MethodGet,
			wantCode:      http.StatusNotFound,
			want:          "",
		},
		{
			name:          "post increment 1",
			postID:        "12345",
			requestMethod: fiber.MethodPost,
			wantCode:      http.StatusCreated,
			want:          "1",
		},
		{
			name:          "post increment 2",
			postID:        "12345",
			requestMethod: fiber.MethodPost,
			wantCode:      http.StatusOK,
			want:          "2",
		},
		{
			name:          "get incremented post",
			postID:        "12345",
			requestMethod: fiber.MethodGet,
			wantCode:      http.StatusOK,
			want:          "2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, tErr := http.NewRequest(
				tc.requestMethod,
				fmt.Sprintf("http://localhost:8080/likes/%s", tc.postID),
				nil,
			)
			tr := require.New(t)
			tr.NoError(tErr)

			httpClient := http.Client{}
			resp, tErr := httpClient.Do(req)
			tr.NoError(tErr)
			defer resp.Body.Close()

			tr.Equal(tc.wantCode, resp.StatusCode)
			if tc.want != "" {
				body, rErr := io.ReadAll(resp.Body)
				tr.NoError(rErr)
				tr.Equal(tc.want, string(body))
			}
		})
	}
}
