package helpers

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"os"
	"whatsgoingon/data"
)

const (
	WhatsAppMessageContentList = 0
)

var (
	redisHostname = os.Getenv("redis_hostname")
	redisPort     = os.Getenv("redis_port")
	redisPassword = os.Getenv("redis_password")
)

func SendMessageToRedis(content data.StoredMessage, clientID string) {
	ctx := context.Background()
	// Create a new Redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Use default Addr
		Password: "admin",
		DB:       WhatsAppMessageContentList, // Use default DB
	})

	// Ping the Redis server and check if any errors occurred.
	err := client.Ping(ctx).Err()
	if err != nil {
		failOnError(err, "Failed to ping Redis server")
		return
	}

	// Transform content to JSON using marshall
	jsonContent, err := json.Marshal(content)
	if err != nil {
		failOnError(err, "Failed to marshal content to JSON")
		return
	}

	// Save the JSON to Redis using the client's Set method.
	err = client.RPush(ctx, clientID, jsonContent).Err()
	if err != nil {
		failOnError(err, "Failed to save JSON to Redis")
		return
	}

	// Close the Redis client.
	err = client.Close()
	if err != nil {
		failOnError(err, "Failed to close Redis client")
		return
	}
	log.Printf("Message sent to Redis: %s", clientID)
}
