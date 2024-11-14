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
