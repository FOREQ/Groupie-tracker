package helpers

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrRootDirNotFound = errors.New("could not find project root directory")

func ChangeDirProjectRoot() error {
	for {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		if _, err := os.Lstat(filepath.Join(dir, "go.mod")); err == nil {
			return nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ErrRootDirNotFound
		}
		if err := os.Chdir(parent); err != nil {
			return err
		}
	}
}
