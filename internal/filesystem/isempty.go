package filesystem

import (
	"errors"
	"io"
	"os"
)

func DirChecks(path string) (exists, empty bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return exists, empty, nil
		}
		return exists, empty, err
	}
	defer f.Close()

	exists = true

	_, err = f.ReadDir(1)
	if errors.Is(err, io.EOF) {
		empty = true
		return exists, empty, nil
	}

	return exists, empty, err
}
