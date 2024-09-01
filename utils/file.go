package utils

import (
	"fmt"
	"os"
)

// Function to delete a file (logo) after use
func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	HandleErr("Error: ", err)
	return nil
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("environment variable %s not set", key)
	}
	return value
}
