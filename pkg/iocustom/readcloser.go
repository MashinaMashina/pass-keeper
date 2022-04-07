package iocustom

import "io"

type ReadCloser struct {
	io.Reader
}

func (receiver ReadCloser) Close() error {
	return nil
}
