package filesystem

import (
	"io/fs"
	"path/filepath"
)

func WalkSrv(path string) (files, directories []string, err error) {
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			directories = append(directories, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	return files, directories, err
}
