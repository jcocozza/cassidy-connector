package utils

import "os"

// write a slice of bytes to a file
func WriteOutput(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}