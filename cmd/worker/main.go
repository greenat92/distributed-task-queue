package main

import (
	"context"
	"distributed-task-queue/internal"
	"distributed-task-queue/internal/monitoring"
	"distributed-task-queue/internal/queue"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	// setup shudown handling
	ctx, cancel := context.WithCancel(context.Background())
	go handleShutdown(cancel)

	log.Println("Starting Worker...")

	// Start polling for tasks
	startWorker(ctx)
}

func processTask(task string) {
	startTime := time.Now()

	parsedTask, err := parseTaskString(task)
	if err != nil {
		log.Printf("Error parsing task: %v", err)
		return
	}

	taskID := parsedTask.ID
	payload := parsedTask.Payload
	statusQueue := "task_status"
	statusInProgress := "in-progress"
	statusFailed := "failed"
	statusCompleted := "completed"

	// Update task status to "in-progress"
	err = redisQueue.SetTaskStatus(statusQueue, taskID, statusInProgress)
	if err != nil {
		log.Printf("Failed to update task status: %v", err)
		return
	}

	log.Printf("Processing task: %s", taskID)

	// Simulate task processing with a delay
	time.Sleep(3 * time.Second)

	// Simulate success or failure
	if strings.Contains(payload, "fail") {
		log.Printf("Task failed: %s", taskID)
		redisQueue.SetTaskStatus(statusQueue, taskID, statusFailed)
		enqueueRetry(task) // Retry the task

		// Update metrics
		monitoring.TasksProcessed.WithLabelValues("failed").Inc()
		monitoring.TaskRetries.Inc()
		return
	}

	log.Printf("Task completed: %s", taskID)
	redisQueue.SetTaskStatus(statusQueue, taskID, statusCompleted)

	// Update metrics
	monitoring.TasksProcessed.WithLabelValues("success").Inc()
	monitoring.TaskProcessingTime.Observe(float64(time.Since(startTime).Seconds()))
}

// enqueue retry implementation
func enqueueRetry(task string) {
	parsedTask, err := parseTaskString(task)
	if err != nil {
		log.Printf("Error parsing task for retry: %v", err)
		return
	}

	taskID := parsedTask.ID
	payload := parsedTask.Payload
	retryQueue := "task_retries"
	taskQueue := "task_queue"
	statusQueue := "task_status"

	statusFailedPermanently := "failed-permanently"
	statusInProgress := "in-progress"

	retryCount, err := redisQueue.IncrementRetryCount(retryQueue, taskID)
	if err != nil {
		log.Printf("Failed to increment retry count for task %s: %v", taskID, err)
		return
	}

	maxRetries := 3
	if retryCount > int64(maxRetries) {
		log.Printf("Task %s exceeded max retries", taskID)
		redisQueue.SetTaskStatus(statusQueue, taskID, statusFailedPermanently)
		return
	}

	log.Printf("Re-enqueuing task %s for retry (%d/%d)", taskID, retryCount, maxRetries)
	err = redisQueue.Enqueue(taskQueue, taskID, payload, statusInProgress)
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
	return internal.Task{
		ID:      parts[0],
		Payload: parts[1],
		Status:  parts[2],
	}, nil
}

func startWorker(ctx context.Context) {
	taskQueue := "task_queue"
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down worker...")
			return
		default:
			// pool for tasks
			task, err := redisQueue.Dequeue(taskQueue)
			if err != nil {
				if strings.Contains(err.Error(), "nil") {
					time.Sleep(1 * time.Second) // No task available, wait before retrying
					continue
				}
				log.Printf("Error dequeuing task: %v", err)
				time.Sleep(2 * time.Second) // Retry after a short delay
				continue
			}
			if task != "" {
				log.Printf("Task dequeued: %s", task)
				go processTask(task) // process each task in a separate goroutine
			}
		}
	}
}

func handleShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	log.Println("Received shutdown signal")
	cancel()
}
