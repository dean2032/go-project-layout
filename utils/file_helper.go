package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// IsFileExist return true if file exist
func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDirectory return true if path is directory
func IsDirectory(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

// EnsureDirExist create directory specified by filePath
func EnsureDirExist(filePath string) error {
	dir := filepath.Dir(filePath)
	isDir, _ := IsDirectory(dir)
	if !isDir {
		err := os.MkdirAll(dir, os.ModeDir|0755)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// IsLink return true if path is a symbol link
func IsLink(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode&os.ModeSymlink != 0:
		return true, nil
	}

	return false, nil
}

// IsDirLink return true if path is a directory and symbol link
func IsDirLink(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode&(os.ModeSticky|os.ModeDir) == (os.ModeSticky | os.ModeDir):
		return true, nil
	}

	return false, nil
}

// IsDirEmpty return true if directory is empty
func IsDirEmpty(path string) bool {
	dir, err := os.ReadDir(path)
	if err != nil {
		return true
	}
	if len(dir) == 0 {
		return true
	}

	return false
}
