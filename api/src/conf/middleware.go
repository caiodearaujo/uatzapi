package conf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"sync"
)

var (
	token      string
	tokenMutex sync.Mutex
)

// Function to initialize the token
func InitToken() {
	token = os.Getenv("API_KEY_TOKEN")
	if token == "" {
		token = uuid.New().String()
		fmt.Println("A new token was generated:", token)
	}
}

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenMutex.Lock()
		defer tokenMutex.Unlock()

		// Verify if request token is valid
		if c.Request.Header.Get("X-Api-Key") != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Continue to the next middleware or handler
		c.Next()
	}
}

func CORSmiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, x-api-key, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
