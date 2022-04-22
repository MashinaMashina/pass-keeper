package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type File interface {
	Name() string
	FullPath() string
	Dir() string
	Size() int64
	Mode() fs.FileMode
	ModTime() time.Time
	IsDir() bool
	Sys() interface{}
	IsHidden() bool
}

type file struct {
	fs.FileInfo
	dir string
}

func NewFile(info fs.FileInfo, dir string) File {
	return file{info, filepath.Clean(dir)}
}

func NewFileFromPath(path string) (File, error) {
	lstat, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	dir, _ := filepath.Split(path)
	return file{lstat, filepath.Clean(dir)}, nil
}

func (f file) FullPath() string {
	return f.dir + string(os.PathSeparator) + f.Name()
}

func (f file) Dir() string {
	return f.dir
}
