package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}


	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Welcome to the Simple Server!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Hello there, web wanderer!")
}

func main() {
	addr := ":8080"
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting simple HTTP server on address %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
