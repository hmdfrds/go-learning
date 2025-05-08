package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TodoItem struct {
	ID          int       `json:"id"`
	Task        string    `json:"task"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
}

var (
	db *sql.DB
)

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error while opening dbPath %s: %v", dbPath, err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatalf("Error while pinging db: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS todos (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	task TEXT NOT NULL,
	description TEXT,
	completed BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		db.Close()
		log.Fatalf("Error while running create table query: %v", err)
	}
	fmt.Println("Database initialized successfully.")
	return db, nil

}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`)) // Simple error message
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	errorPayload := map[string]string{"error": message}
	respondWithJSON(w, statusCode, errorPayload)
}

type CreateTodoRequest struct {
	Task        string `json:"task"`
	Description string `json:"description"`
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			respondWithError(w, http.StatusBadRequest, msg)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			respondWithError(w, http.StatusBadRequest, msg)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			respondWithError(w, http.StatusBadRequest, msg)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			respondWithError(w, http.StatusBadRequest, msg)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			respondWithError(w, http.StatusBadRequest, msg)
		default:
			log.Println(err.Error()) // Log other unexpected errors
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	if req.Task == "" {
		respondWithError(w, http.StatusBadRequest, "Field 'task' is required.")
		return
	}

	newItem := TodoItem{
		Task:        req.Task,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
	}
	query := "INSERT INTO todos(task, description, completed, created_at) VALUES(?, ?, ?, ?) RETURNING id;"
	err = db.QueryRowContext(r.Context(), query, newItem.Task, newItem.Description, newItem.Completed, newItem.CreatedAt).Scan(&newItem.ID)
	if err != nil {
		respondWithError(w, http.StatusBadGateway, fmt.Sprintf("Error insert todo: %v", err))
	}

	log.Printf("Created Todo Item: ID=%d, Task=%s\n", newItem.ID, newItem.Task)
	respondWithJSON(w, http.StatusCreated, newItem)
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {

	query := "SELECT id, task, description, completed, created_at FROM todos ORDER BY created_at DESC"
	rows, err := db.QueryContext(r.Context(), query)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error while selecting todo: %v", err))
	}
	defer rows.Close()

	var itemsToReturn []TodoItem
	for rows.Next() {
		var todo TodoItem
		if err := rows.Scan(&todo.ID, &todo.Task, &todo.Description, &todo.Completed, &todo.CreatedAt); err != nil {
			log.Print("Error while parsing row")
		}
		itemsToReturn = append(itemsToReturn, todo)
	}

	if rows.Err() != nil {
		respondWithError(w, http.StatusBadGateway, fmt.Sprintf("Error while parsing row: %v", err))
	}

	log.Printf("Retrieved %d Todo Items\n", len(itemsToReturn))
	respondWithJSON(w, http.StatusOK, itemsToReturn)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTodos(w, r)
	case http.MethodPost:
		handleCreateTodo(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method %s not allowed for this endpoint.", r.Method))
	}
}

func main() {
	var err error
	db, err = initDB("todos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/todos", todosHandler)

	port := ":8080"
	log.Printf("Server starting on http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
