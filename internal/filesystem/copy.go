package filesystem

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CopyDir(sourcePath, destinationPath string) (err error) {
	return filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		targetPath := filepath.Join(destinationPath, strings.TrimPrefix(path, sourcePath))
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0700)
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, b, 0400)
	})
}
