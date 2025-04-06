package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	testCases := []struct {
		name     string
		courseID string
		want     string
	}{
		{
			name:     "course 1",
			courseID: "1",
			want:     "Introduction to programming",
		},
		{
			name:     "course 2",
			courseID: "2",
			want:     "Introduction to algorithms",
		},
		{
			name:     "course 3",
			courseID: "3",
			want:     "Data structures",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpRequest, err := http.NewRequest("GET", "/courses/description?course_id="+tc.courseID, nil)
			r := require.New(t)
			r.NoError(err)

			httpRecorder := httptest.NewRecorder()

			handler := http.HandlerFunc(CourseDescHandler)
			handler.ServeHTTP(httpRecorder, httpRequest)

			r.Equal(http.StatusOK, httpRecorder.Code)
			r.Equal(tc.want, httpRecorder.Body.String())
		})
	}
}
