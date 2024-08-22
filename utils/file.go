package utils

import (
	"fmt"
	"os"
)

// Function to delete a file (logo) after use
func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	fmt.Println("Deleted file:", filepath)
	return nil
}
