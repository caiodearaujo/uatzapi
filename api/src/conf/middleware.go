package conf

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	token      string
	tokenMutex sync.RWMutex // Use RWMutex for better performance in concurrent read scenarios.
)

// InitToken initializes the authentication token from environment variables.
// If the token is not set, a new UUID token is generated.
func InitToken() {
	tokenMutex.Lock() // Lock the token for write access.
	defer tokenMutex.Unlock()

	// Retrieve the token from the environment variable.
	token = os.Getenv("API_KEY_TOKEN")
	if token == "" {
		// If no token is found in the environment, generate a new UUID.
		token = uuid.New().String()
		fmt.Println("A new token was generated:", token)
	} else {
		fmt.Println("Token loaded from environment:", token)
	}
}

// TokenMiddleware is a Gin middleware that validates incoming requests using the `X-Api-Key` header.
// It checks if the provided token matches the one initialized by InitToken.
// Unauthorized requests receive a 401 response.
func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lock the token for read access with a read lock.
		tokenMutex.RLock()
		defer tokenMutex.RUnlock()

		// Check if the request contains the correct API token.
		if c.Request.Header.Get("X-Api-Key") != token {
			// Respond with 401 Unauthorized if the token is invalid.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Proceed to the next middleware or handler if the token is valid.
		c.Next()
	}
}

// CORSMiddleware is a Gin middleware that handles Cross-Origin Resource Sharing (CORS) headers.
// It allows requests from any origin and supports credentials, specific headers, and HTTP methods.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers to allow requests from any origin.
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-Api-Key, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		// Respond to preflight OPTIONS requests and terminate them early.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Proceed to the next middleware or handler if the request method is not OPTIONS.
		c.Next()
	}
}
