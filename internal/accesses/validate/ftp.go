package validate

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"pass-keeper/internal/accesses/accesstype"
	"time"
)

func FTP(access accesstype.Access, verifyUnknownHosts bool, w io.Writer, r io.Reader) (bool, error) {
	addr := fmt.Sprintf("%s:%d", access.Host(), access.Port())

	conn, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return false, err
	}

	defer conn.Quit()

	err = conn.Login(access.Login(), access.Password())
	if err != nil {
		return false, err
	}

	return true, nil
}
