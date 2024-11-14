package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"distributed-task-queue/internal"
	"distributed-task-queue/internal/queue"

	"github.com/gorilla/mux"
)

var redisQueue *queue.RedisQueue

func main() {
	// init the redis queue
	redisQueue = queue.NewRedisQueue("localhost:6379")

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
	task.ID = "task" + fmt.Sprint(time.Now().UnixNano())
	task.Status = "submitted"

	// Engueue the task into Redis queue
	if err := redisQueue.Enqueue("task_queue", task.ID, task.Payload, task.Status); err != nil {
		http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
		return
	}

	// Respond to the client with the task ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID})
}
