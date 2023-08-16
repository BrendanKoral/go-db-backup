package fs

import (
	"errors"
	"os"
)

// CreateBackupDir creates the backup directory if it doesn't exist
func CreateBackupDir(path string) error {
	var err error

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	// This is ok, dir already exists
	if _, err := os.Stat(path); errors.Is(err, os.ErrExist) {
	}

	return err
}
