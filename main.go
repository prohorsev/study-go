package main

import (
	"net/http"
	"strconv"
	"time"
)

var courses = map[int64]string{
	1: "Introduction to programming",
	2: "Introduction to algorithms",
	3: "Data structures",
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/courses/description", CourseDescHandler)

	server := &http.Server{
		Addr:              ":80",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Go to /courses/description"))
}

func CourseDescHandler(w http.ResponseWriter, r *http.Request) {
	cId := r.URL.Query().Get("course_id")
	id, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(courses[id]))
}
