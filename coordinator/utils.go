package coordinator

import (
	"os"
)

// CreateFile creates a file at the given path and returns the file handle.
// It returns an error if the file could not be created.
func CreateFile(path string) (*os.File, error) {
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return out, nil
}
