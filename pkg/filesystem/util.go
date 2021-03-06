package filesystem

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Stat(path string) (File, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	dirname, _ := filepath.Split(path)

	return NewFile(fileInfo, dirname), nil
}

func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func ReadDir(dirname string) ([]File, error) {
	list, err := ioutil.ReadDir(dirname)

	if err != nil {
		return nil, err
	}

	files := make([]File, 0)
	for _, fsFile := range list {
		files = append(files, NewFile(fsFile, dirname))
	}

	return files, nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func PrettyPath(s string) string {
	return strings.ReplaceAll(s, "\\", "/")
}
