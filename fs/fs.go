package fs

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DumpDir string
	Test    string
}

var timeFormat = "20060102T150405"

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

func DeleteInvalidBackupFile(path string, fileName string) error {
	file := fmt.Sprintf("%s/%s", path, fileName)

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return err
	}

	err := os.Remove(file)

	return err
}
