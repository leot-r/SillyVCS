package utils

import (
	"os"
)

func CheckIfInit(dotfilePath string) (exists bool, error error) {
	// Check if a repo has been initalized in the current directory
	_, err := os.Stat(dotfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

