package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	testCases := []struct {
		name     string
		from     string
		to       string
		wantCode int
		want     string
	}{
		{
			name:     "exchange rate not found 1",
			from:     "USD",
			to:       "GEL",
			wantCode: http.StatusNotFound,
			want:     "",
		},
		{
			name:     "exchange rate not found 2",
			from:     "GEL",
			to:       "USD",
			wantCode: http.StatusNotFound,
			want:     "",
		},
		{
			name:     "positive 1",
			from:     "EUR",
			to:       "USD",
			wantCode: http.StatusOK,
			want:     "1.25",
		},
		{
			name:     "positive 2",
			from:     "USD",
			to:       "JPY",
			wantCode: http.StatusOK,
			want:     "110.00",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, tErr := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("http://localhost:8080/convert?from=%s&to=%s", tc.from, tc.to),
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
