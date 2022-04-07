package scanln

import (
	"bufio"
	"io"
	"os"
)

func Scanln(output *string) error {
	s := bufio.NewScanner(os.Stdin)
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
