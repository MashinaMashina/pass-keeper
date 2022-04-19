package validate

import (
	"io"
	"pass-keeper/internal/accesses/accesstype"
)

var Validators = map[string]func(access accesstype.Access, verifyUnknownHosts bool, w io.Writer, r io.Reader) (bool, error){
	"ssh": SSH,
	"ftp": FTP,
}
