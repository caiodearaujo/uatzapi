package handler

import (
	"log"
)

// FailOnError logs an error message if the given error is not nil.
// This function is useful for handling and logging errors without panicking the application.
// If an error is encountered, it logs the custom message along with the error details.
func FailOnError(err error, msg string) {
	if err != nil {
		// Log the error with the custom message and the error details.
		log.Printf("[ERROR] %s: %v", msg, err)
	}
}
