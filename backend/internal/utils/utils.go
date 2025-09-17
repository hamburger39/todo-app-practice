package utils

import "time"

// GenerateID generates a unique ID based on current timestamp
func GenerateID() string {
	return time.Now().Format("20060102150405")
}
