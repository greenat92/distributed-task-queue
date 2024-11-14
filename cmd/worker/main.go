package main

import (
	"distributed-task-queue/internal"
	"distributed-task-queue/internal/queue"
	"fmt"
	"log"
	"strings"
	"time"
)

var redisQueue *queue.RedisQueue

func main() {
	// Initialize Redis connection
	redisQueue = queue.NewRedisQueue("localhost:6379")

	log.Println("Starting Worker...")

	// Start polling for tasks
	for {
		task, err := redisQueue.Dequeue("task_queue")
		if err != nil {
			log.Printf("Error dequeuing task: %v", err)
			time.Sleep(2 * time.Second) // Retry after a short delay
			continue
		}

		if task != "" {
			log.Printf("Task dequeued: %s", task)
			go processTask(task) // Process each task in a separate goroutine
		}

		time.Sleep(1 * time.Second) // Polling interval
	}
}

func processTask(task string) {
	log.Printf("Processing task: %s", task)
	// Simulate task processing with a delay
	time.Sleep(3 * time.Second)
	log.Printf("Task completed: %s", task)
	// after proceesing dequeue the tasks again
	status := "completed"
	parsedTask, err := parseTaskString(task)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Engueue the task into Redis queue
	if err := redisQueue.Enqueue("completed_task_queue", parsedTask.ID, parsedTask.Payload, status); err != nil {
		log.Printf("something happen while enqueuing the completed task %s", err)
	}
	log.Printf("task enqueued successfuly %s", parsedTask.ID+"|"+parsedTask.Payload+"|"+status)
}

// parseTaskString parses the input string and returns a Task struct.
func parseTaskString(input string) (internal.Task, error) {
	parts := strings.Split(input, "|")
	if len(parts) < 2 {
		return internal.Task{}, fmt.Errorf("input does not contain enough elements separated by '|'")
	}

	// Create a Task with ID and Status from the first two elements
	task := internal.Task{
		ID:      parts[0],
		Payload: parts[1],
		Status:  parts[2],
	}

	return task, nil
}
