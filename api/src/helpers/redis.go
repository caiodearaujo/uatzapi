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
	// WhatsAppMessageContentList represents the Redis database index to use.
	WhatsAppMessageContentList = 0
)

var (
	redisClientInstance *redis.Client
	redisOnce           sync.Once
)

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, the function returns the fallback value provided.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getRedisClient initializes and returns a Redis client using the singleton pattern.
// The client is only created once during the application's lifetime.
func getRedisClient() *redis.Client {
	redisHostname := getEnv("redis_hostname", "localhost")
	redisPort := getEnv("redis_port", "6379")
	redisPassword := getEnv("redis_password", "admin")

	// Ensure the Redis client is only initialized once.
	redisOnce.Do(func() {
		log.Info().Msgf("Creating a new Redis client on dsn: %s:%s", redisHostname, redisPort)
		redisClientInstance = redis.NewClient(&redis.Options{
			Addr:     redisHostname + ":" + redisPort,
			Password: redisPassword,
			DB:       WhatsAppMessageContentList, // Use the specific Redis DB index for message content.
		})
	})

	return redisClientInstance
}

// PingRedis checks the connectivity to the Redis server by sending a ping request.
// It returns an error if the Redis server is unreachable.
func PingRedis(ctx context.Context) error {
	client := getRedisClient()
	return client.Ping(ctx).Err() // Send a ping to Redis and return any errors.
}

// SendMessageToRedis pushes a message to a Redis list for a given device.
// The message is first marshaled to JSON before being sent.
func SendMessageToRedis(ctx context.Context, content data.StoredMessage, deviceID int) {
	// Ping Redis to ensure it's reachable.
	if err := PingRedis(ctx); err != nil {
		handler.FailOnError(err, "Failed to ping Redis server")
		return
	}

	// Get the Redis client.
	client := getRedisClient()

	// Marshal the message content to JSON format.
	jsonContent, err := MarshalMessageToJSON(content)
	if err != nil {
		handler.FailOnError(err, "Failed to marshal content to JSON")
		return
	}

	// Push the JSON content to a Redis list where the key is the device ID.
	if err := client.RPush(ctx, strconv.Itoa(deviceID), jsonContent).Err(); err != nil {
		handler.FailOnError(err, "Failed to push message to Redis")
		return
	}

	// Log the success of the operation.
	log.Printf("Message sent to Redis successfully for deviceID: %d", deviceID)
}
