package main_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	SignInResponse struct {
		JWTToken string `json:"jwt_token"`
	}

	ProfileResponse struct {
		Email string `json:"email"`
	}
)

func TestPractice(t *testing.T) {
	r := require.New(t)

	// Sign up a user
	testRequest(
		r,
		http.MethodPost,
		"/signup",
		"",
		`{"email":"test@test.com","password":"qwerty"}`,
		http.StatusOK,
		nil,
	)

	// Try to sign in with already existed email
	testRequest(
		r,
		http.MethodPost,
		"/signup",
		"",
		`{"email":"test@test.com","password":"foobar"}`,
		http.StatusConflict,
		nil,
	)

	// Try to sign in with wrong email
	testRequest(
		r,
		http.MethodPost,
		"/signin",
		"",
		`{"email":"test2@test.com","password":"qwerty"}`,
		http.StatusUnprocessableEntity,
		nil,
	)

	// Try to sign in with wrong password
	testRequest(
		r,
		http.MethodPost,
		"/signin",
		"",
		`{"email":"test@test.com","password":"qwerty123"}`,
		http.StatusUnprocessableEntity,
		nil,
	)

	// Sign in with the user
	resp := SignInResponse{}
	testRequest(
		r,
		http.MethodPost,
		"/signin",
		"",
		`{"email":"test@test.com","password":"qwerty"}`,
		http.StatusOK,
		&resp,
	)

	// Get profile of the user
	profileResp := ProfileResponse{}
	testRequest(
		r,
		http.MethodGet,
		"/profile",
		resp.JWTToken,
		"",
		http.StatusOK,
		&profileResp,
	)
	r.Equal("test@test.com", profileResp.Email)

	// Try to get profile with invalid JWT token
	testRequest(
		r,
		http.MethodGet,
		"/profile",
		"invalid",
		"",
		http.StatusUnauthorized,
		nil,
	)
}

func testRequest(r *require.Assertions, method, path, jwtToken string, body string, wantCode int, responseBody interface{}) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, "http://localhost:8080"+path, bodyReader)
	r.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	if jwtToken != "" {
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	r.NoError(err)

	r.Equal(wantCode, resp.StatusCode)

	if responseBody != nil {
		bodyBytes, jErr := io.ReadAll(resp.Body)
		r.NoError(jErr)

		jErr = json.Unmarshal(bodyBytes, responseBody)
		r.NoError(jErr)
	}
}
