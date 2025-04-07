package main

import (
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, ".log")
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	logger.SetOutput(file)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go to /sum"))
	})

	mux.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		x, err := strconv.Atoi(r.URL.Query().Get("x"))
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			panic(err)
		}
		if (y > 0 && x > math.MaxInt-y) || (y < 0 && x < math.MinInt-y) {
			logger.WithFields(logrus.Fields{
				"x": x,
				"y": y,
			}).Warning("Sum overflows int")
			w.Write([]byte("-1"))
			return
		}
		sum := x + y
		w.Write([]byte(strconv.Itoa(sum)))
	})

	port := "8080"
	logWithPort := logrus.WithFields(logrus.Fields{
		"port": port,
	})
	logWithPort.Info("Starting a web-server on port")
	logWithPort.Fatal(server.ListenAndServe())
}
