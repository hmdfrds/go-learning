package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ServerInfo struct {
	ServerName  string    `json:"serverName"`
	Version     string    `json:"version"`
	CurrentTime time.Time `json:"currentTime"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	serverInfo := ServerInfo{ServerName: "MySimpleServer", Version: "1.0.1", CurrentTime: time.Now()}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(serverInfo)
	if err != nil {
		log.Printf("Error encoding /api/info response: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Welcome to the Simple Server!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World"
		}
		response := map[string]string{"greeting": fmt.Sprintf("Hello, %s! (From GET)", name)}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Error encoding /name response: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	case http.MethodPost:
		type HelloRequest struct {
			Name string `json:"name"`
		}
		defer r.Body.Close()
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var req HelloRequest
		err := decoder.Decode(&req)

		if err != nil {
			http.Error(w, "Bad request: Invalid JSON format or data", http.StatusBadRequest)
			return
		}

		name := req.Name
		if name == "" {
			name = "World"
		}

		response := map[string]string{"receivedName": name, "message": fmt.Sprintf("Received POST for %s!", name)}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("Error encoding /name response: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "Hello, %s! (from POST)\n", name)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	addr := ":8080"
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/api/info", infoHandler)

	fmt.Printf("Starting simple HTTP server on address %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}
