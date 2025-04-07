package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	r := require.New(t)
	// Check that we can't run a new process on the same port.
	runOutput2 := bytes.NewBuffer(nil)
	runCmd2 := exec.Command("go", "run", "main.go")
	runCmd2.Stdout = runOutput2
	runCmd2.Stderr = runOutput2
	r.Error(runCmd2.Run())
	r.Contains(runOutput2.String(), "listen tcp :8080: bind: address already in use")

	// Run test cases.
	testCases := []struct {
		name                 string
		x                    int
		y                    int
		expectedResponseBody string
		expectedCode         int
		expectedLog          string
	}{
		{
			name:                 "positive",
			x:                    11,
			y:                    55,
			expectedResponseBody: "66",
			expectedCode:         http.StatusOK,
			expectedLog:          "",
		},
		{
			name:                 "positive #2",
			x:                    124,
			y:                    77777,
			expectedResponseBody: "77901",
			expectedCode:         http.StatusOK,
			expectedLog:          "",
		},
		{
			name:                 "overflow with bigger x",
			x:                    math.MaxInt - 10,
			y:                    700,
			expectedResponseBody: "-1",
			expectedCode:         http.StatusOK,
			expectedLog:          fmt.Sprintf("level=warning msg=\"Sum overflows int\" x=%d y=700", math.MaxInt-10),
		},
		{
			name:                 "overflow with bigger x",
			x:                    500,
			y:                    math.MaxInt - 10,
			expectedResponseBody: "-1",
			expectedCode:         http.StatusOK,
			expectedLog:          fmt.Sprintf("level=warning msg=\"Sum overflows int\" x=500 y=%d", math.MaxInt-10),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, tErr := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("http://localhost:8080/sum?x=%d&y=%d", tc.x, tc.y),
				nil,
			)
			tr := require.New(t)
			tr.NoError(tErr)

			httpClient := http.Client{}
			resp, tErr := httpClient.Do(req)
			tr.NoError(tErr)
			defer resp.Body.Close()

			tr.Equal(tc.expectedCode, resp.StatusCode)
			if tc.expectedResponseBody != "" {
				body, rErr := io.ReadAll(resp.Body)
				tr.NoError(rErr)
				tr.Equal(tc.expectedResponseBody, string(body))
			}

			dat, _ := os.ReadFile(".log")
			log := string(dat)
			fmt.Print(log)
			if tc.expectedLog != "" {
				tr.Contains(log, tc.expectedLog)
			}
		})
	}
}
