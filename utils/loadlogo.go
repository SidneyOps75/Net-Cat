package utils

import (
	"fmt"
	"os"
)

// LoadLogo reads the contents of the "logo.txt" file and returns it as a string
func LoadLogo() (string, error) {
	content, err := os.ReadFile("logo.txt")
	if err != nil {
		return "", fmt.Errorf("error reading logo file: %v", err)
	}
	return string(content), nil
}
