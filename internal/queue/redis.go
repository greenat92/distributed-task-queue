package queue

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisQueue struct {
	client *redis.Client
}

// NewRedisQue initializes a new Redis client
func NewRedisQueue(addr string) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to redis: %v", err)
	}

	return &RedisQueue{client: client}
}

// engqueue function
func (q *RedisQueue) Enqueue(key string, taskID string, payload string, status string) error {
	// push the task into the Redis list
	err := q.client.RPush(ctx, key, taskID+"|"+payload+"|"+status).Err()
	if err != nil {
		return err
	}
	log.Printf("task enqueued: %s\n", taskID)
	return nil
}

func (q *RedisQueue) Dequeue(queueName string) (string, error) {
	task, err := q.client.LPop(ctx, queueName).Result()
	if err == redis.Nil {
		return "", nil // No task in the queue
	}
	return task, err
}

// SetTaskStatus updates the status of a task in redis
func (q *RedisQueue) SetTaskStatus(key string, taskID string, status string) error {
	err := q.client.HSet(ctx, key, taskID, status).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetTaskStatus retrieves the status of a task from redis
func (q *RedisQueue) GetTaskStatus(key string, taskID string) (string, error) {
	status, err := q.client.HGet(ctx, key, taskID).Result()
	if err == redis.Nil {
		return "", nil // task not found
	}
	return status, err
}

// retry func
func (q *RedisQueue) IncrementRetryCount(key string, taskID string) (int64, error) {
	retryCount, err := q.client.HIncrBy(ctx, key, taskID, 1).Result()

	if err != nil {
		return 0, err
	}
	return int64(retryCount), nil
}
