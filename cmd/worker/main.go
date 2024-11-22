package main

import (
	"distributed-task-queue/internal"
	"distributed-task-queue/internal/monitoring"
	"distributed-task-queue/internal/queue"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var redisQueue *queue.RedisQueue

func main() {

	// init metrics
	monitoring.InitMetrics()

	// serve metrics on a separate goroutine
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics endpoint available at :2112/metrics")
		log.Fatal(http.ListenAndServe("127.0.0.1:2112", nil))
	}()

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
	startTime := time.Now()

	parsedTask, err := parseTaskString(task)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	taskID := parsedTask.ID
	payload := parsedTask.Payload
	queue := "task_status"
	statusInprogress := "in-progress"
	statusFailed := "failed"
	statusCompleted := "completed"

	// Upate status to in porgress
	err = redisQueue.SetTaskStatus(queue, taskID, statusInprogress)
	if err != nil {
		log.Printf("Failed to update task status: %v", err)
		return
	}

	log.Printf("processing task: %s", taskID)

	// simulate task processing with a delay
	time.Sleep(3 * time.Second)

	// simulate processing success or failure
	if strings.Contains(payload, "fail") {
		log.Printf("Task Failed: %s", taskID)
		redisQueue.SetTaskStatus(queue, taskID, statusFailed)
		enqueueRetry(task) // retry the task

		// update metrics
		monitoring.TasksProcessed.WithLabelValues("failed").Inc()
		monitoring.TaskReties.Inc()
		return
	}

	log.Printf("Task completed: %s", taskID)
	redisQueue.SetTaskStatus(queue, taskID, statusCompleted)

	// update metrics
	monitoring.TasksProcessed.WithLabelValues("success").Inc()
	monitoring.TaskProcessingTime.Observe(float64(time.Since(startTime).Seconds()))
}

// enqueue retry implementation
func enqueueRetry(task string) {
	parsedTask, err := parseTaskString(task)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	taskID := parsedTask.ID
	payload := parsedTask.Payload
	queueRetries := "task_retries"
	queueTask := "task_queue"
	queueTaskStatus := "task_status"

	statusFailedPermanenty := "failed-permanently"
	statusInprogress := "in-progress"

	retryCount, err := redisQueue.IncrementRetryCount(queueRetries, taskID)
	if err != nil {
		log.Printf("Failed to increment retry count for task %s: %V", taskID, err)
		return
	}

	maxRetries := 3
	if retryCount > int64(maxRetries) {
		log.Printf("Task %s exceeded max retries", taskID)
		redisQueue.SetTaskStatus(queueTaskStatus, taskID, statusFailedPermanenty)
		return
	}

	log.Printf("Re-engueuing task %s for retry (%d/%d)", taskID, retryCount, maxRetries)
	err = redisQueue.Enqueue(queueTask, taskID, payload, statusInprogress)
	if err != nil {
		log.Printf("Failed to re-enqueue task %s: %v", taskID, err)
	}

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
