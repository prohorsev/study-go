package main_test

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	r := require.New(t)

	// Send 4 parallel requests: 2 to /foo and 2 to /bar
	// 1 request to /foo and 1 request to bar should be successful
	// other requests should be rejected with 429 status code

	requests := []*http.Request{
		request(r, "/foo"),
		request(r, "/bar"),
		request(r, "/foo"),
		request(r, "/bar"),
	}

	wg := sync.WaitGroup{}

	for _, req := range requests {
		wg.Add(1)
		go func(req *http.Request) {
			defer wg.Done()

			_, gErr := http.DefaultClient.Do(req)
			r.NoError(gErr)
		}(req)
	}

	wg.Wait()

	data, _ := os.ReadFile(".log")
	output := string(data)

	lines := strings.Split(output, "\n")

	expectedOutputs := []string{
		": GET /foo - 200",
		": GET /bar - 200",
		": GET /foo - 429",
		": GET /bar - 429",
	}

	for _, expectedOutput := range expectedOutputs {
		r.Contains(output, expectedOutput)

		for _, line := range lines {
			if !strings.HasPrefix(line, expectedOutput) {
				continue
			}

			requestID := strings.TrimSuffix(line, expectedOutput)
			r.True(IsValidUUID(requestID), "Invalid request ID in line: %s", line)
			break
		}
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func request(r *require.Assertions, path string) *http.Request {
	req, tErr := http.NewRequest(http.MethodGet, "http://localhost:8080"+path, nil)
	r.NoError(tErr)

	return req
}
