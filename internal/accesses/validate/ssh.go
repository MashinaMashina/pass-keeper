package validate

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"pass-keeper/internal/accesses/accesstype"
	"strings"
)

func SSH(access accesstype.Access, verifyUnknownHosts bool, w io.Writer, r io.Reader) (bool, error) {
	var auth goph.Auth

	if access.Password() != "" {
		auth = goph.Password(access.Password())
	}

	client, err := goph.NewConn(&goph.Config{
		User: access.Login(),
		Addr: access.Host(),
		Port: uint(access.Port()),
		Auth: auth,
		Callback: func(host string, remote net.Addr, key ssh.PublicKey) error {
			if !verifyUnknownHosts {
				return nil
			}

			hostFound, _ := goph.CheckKnownHost(host, remote, key, "")

			if !hostFound {
				fmt.Fprintln(w, "Unknown host", host, remote)
				fmt.Fprintln(w, key.Type()+" "+base64.StdEncoding.EncodeToString(key.Marshal()))
				fmt.Fprint(w, "Trust this host? (Y/n) ")

				var val string
				fmt.Fscanln(r, &val)

				if strings.ToLower(val) == "y" {
					return goph.AddKnownHost(host, remote, key, "")
				}

				return errors.New("invalid host")
			}

			return nil
		},
	})
	if err != nil {
		return false, err
	}
	defer client.Close()

	out, err := client.Run("ls -a ~")
	if err != nil {
		return false, err
	}

	if out == nil || len(out) == 0 {
		return false, errors.New("missed output")
	}

	return true, nil
}
