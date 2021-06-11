package filesystem

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed dummy/*
var dummySrv embed.FS //nolint:gochecknoglobals

func InitDummySrv(targetRootPath string) (err error) {
	const pathPrefix = "dummy"
	return fs.WalkDir(dummySrv, pathPrefix, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		targetPath := filepath.Join(targetRootPath, strings.TrimPrefix(path, pathPrefix))
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0400)
		}

		b, err := dummySrv.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, b, 0400)
	})
}
