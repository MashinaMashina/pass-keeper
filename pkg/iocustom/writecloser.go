package iocustom

import "io"

type WriteCloser struct {
	io.Writer
}

func (receiver WriteCloser) Close() error {
	return nil
}
