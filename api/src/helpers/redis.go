package helpers

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"

	"whatsgoingon/data"
	"whatsgoingon/handler"
)

const (
	WhatsAppMessageContentList = 0
)

var (
	redisClientInstance *redis.Client
	redisOnce           sync.Once
)

// Get environment variable with a fallback value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Initialize a Redis client once using the sync.Once pattern to ensure a singleton instance.
func getRedisClient() *redis.Client {
	redisHostname := os.Getenv("REDIS_HOSTNAME")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisOnce.Do(func() {
		log.Info().Msg("Creating a new Redis client on dsn: " + redisHostname + ":" + redisPort)
		redisClientInstance = redis.NewClient(&redis.Options{
			Addr:     redisHostname + ":" + redisPort,
			Password: redisPassword,
			DB:       WhatsAppMessageContentList,
		})
	})

	return redisClientInstance
}

// PingRedis checks if Redis is available by pinging the server.
func PingRedis(ctx context.Context) error {
	client := getRedisClient()
	return client.Ping(ctx).Err()
}

func SendMessageToRedis(ctx context.Context, content data.StoredMessage, deviceID int) {
	// Ping the Redis server and check if any errors occurred.
	if err := PingRedis(ctx); err != nil {
		handler.FailOnError(err, "Failed to ping Redis server")
		return
	}

	// Create a new Redis client.
	client := getRedisClient()

	// Marshall content to JSON
	jsonContent, err := MarshalMessageToJSON(content)
	if err != nil {
		handler.FailOnError(err, "Failed to marshal content to JSON")
		return
	}

	// Save the JSON to Redis using the client's Set method.
	if err := client.RPush(ctx, strconv.Itoa(deviceID), jsonContent).Err(); err != nil {
		handler.FailOnError(err, "Failed to save JSON to Redis")
		return
	}

	log.Printf("Message sent to Redis successfully for deviceID: %s", deviceID)
}
