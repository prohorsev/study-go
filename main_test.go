package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	BinarySearchResponse struct {
		TargetIndex int    `json:"target_index"`
		Error       string `json:"error,omitempty"`
	}
)

func TestPractice(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
		wantCode    int
		wantIndex   int
		wantError   string
	}{
		{
			name:        "invalid JSON request body",
			requestBody: `{"numbe`,
			wantCode:    http.StatusBadRequest,
			wantIndex:   -1,
			wantError:   "Invalid JSON",
		},
		{
			name:        "target number is not found",
			requestBody: `{"numbers": [1, 2, 3, 5], "target": 4}`,
			wantCode:    http.StatusNotFound,
			wantIndex:   -1,
			wantError:   "Target was not found",
		},
		{
			name:        "numbers count is even. Target number is found",
			requestBody: `{"numbers": [1, 2, 3, 7, 99, 100, 250, 1000], "target": 99}`,
			wantCode:    http.StatusOK,
			wantIndex:   4,
		},
		{
			name:        "numbers count is odd. Target number is found",
			requestBody: `{"numbers": [1, 2, 7, 99, 100, 250, 1000], "target": 99}`,
			wantCode:    http.StatusOK,
			wantIndex:   3,
		},
	}

	runOutput := bytes.NewBuffer(nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer fmt.Printf("Web-app output: %s\n", runOutput.String())

			req, tErr := http.NewRequest(
				http.MethodPost,
				"http://localhost:8080/search",
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

			bsResp := BinarySearchResponse{}
			tErr = json.Unmarshal(bodyBytes, &bsResp)
			tr.NoError(tErr)

			tr.Equal(tc.wantIndex, bsResp.TargetIndex)
			tr.Equal(tc.wantError, bsResp.Error)
		})
	}
}
