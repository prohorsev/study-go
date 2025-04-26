package main_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type slowConn struct {
	net.Conn
	sr slowReader
}

func newSlowConn(conn net.Conn) net.Conn {
	return slowConn{conn, slowReader{conn}}
}

func (conn slowConn) Read(data []byte) (int, error) {
	return conn.sr.Read(data)
}

type slowReader struct {
	r io.Reader
}

func (r slowReader) Read(data []byte) (int, error) {
	// wait for 500 ms before reading a single byte.
	time.Sleep(500 * time.Millisecond)
	n, err := r.r.Read(data[:1])
	if n > 0 {
		fmt.Printf("%s", data[:1])
	}
	return n, err
}

func TestPractice(t *testing.T) {

	t.Run("bad_request", func(t *testing.T) {
		tr := require.New(t)

		testRequest(tr, "{", fiber.StatusBadRequest, "Invalid JSON")
	})

	t.Run("panics", func(t *testing.T) {
		tr := require.New(t)

		testRequest(tr, `{"user_id":1,"message":"Hello world!"}`, fiber.StatusOK, "OK")
		testRequest(tr, `{"user_id":2,"message":"Hello world2!"}`, fiber.StatusOK, "OK")
		testRequest(tr, `{"user_id":3,"message":"Hello world3!"}`, fiber.StatusOK, "OK")
		testRequest(tr, `{"user_id":4,"message":"Hello world4!"}`, fiber.StatusInternalServerError, "Queue is full")
	})
}

func testRequest(r *require.Assertions, body string, wantCode int, wantBody string) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/push/send", bodyReader)
	r.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	r.NoError(err)

	bodyBytes, err := io.ReadAll(resp.Body)
	r.NoError(err)

	r.Equal(wantCode, resp.StatusCode)
	r.Equal(wantBody, string(bodyBytes))
}
