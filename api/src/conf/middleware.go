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
