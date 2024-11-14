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
func (q *RedisQueue) Enqueue(taskID string, payload string) error {
	// push the task into the Redis list
	err := q.client.RPush(ctx, "task_queue", taskID+"|"+payload).Err()
	if err != nil {
		return err
	}
	log.Printf("task enqueued: %s\n", taskID)
	return nil
}
