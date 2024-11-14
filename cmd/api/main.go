package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"distributed-task-queue/internal"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", handleTaskSubmission).Methods("POST")

	fmt.Println("Starting API server on: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleTaskSubmission(w http.ResponseWriter, r *http.Request) {
	var task internal.Task

	// Pasre json request body
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// basic validation: check if Payload is provided
	if task.Payload == "" {
		http.Error(w, "Payload is required", http.StatusBadRequest)
		return
	}

	// generate a simple ID for the task (e.g., using current timestamp)
	task.ID = "tas-" + fmt.Sprint(time.Now().UnixNano())
	task.Status = "submitted"

	// for now, just log the task as a placeholder for queueing logic
	log.Printf("Task received: %+v\n", task)

	// Respond to the client with the task ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID})
}
