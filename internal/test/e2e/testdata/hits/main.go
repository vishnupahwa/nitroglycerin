package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

func main() {
	var count int64
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt64(&count, 1)
		log.Printf("Incremented: count = %d\n", count)
		w.WriteHeader(200)
	})
	http.HandleFunc("/count", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, strconv.FormatInt(count, 10))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
