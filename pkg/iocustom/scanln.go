package iocustom

import (
	"bufio"
	"io"
	"os"
)

func Scanln(output *string) error {
	return Fscanln(os.Stdin, output)
}

func Fscanln(r io.Reader, output *string) error {
	s := bufio.NewScanner(r)
	if s.Scan() {
		*output = s.Text()
		return nil
	}

	err := s.Err()
	if err == nil {
		err = io.EOF
	}

	return err
}
