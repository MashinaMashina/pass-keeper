package parselnk

import (
	"flag"
	"fmt"
	lnk2 "github.com/parsiya/golnk"
	"github.com/pkg/errors"
	"net/url"
	"pass-keeper/internal/accesses"
	"pass-keeper/pkg/filesystem"
	"strings"
)

func sshAccessByLnk(file filesystem.File) (accesses.Access, error) {
	lnk, err := lnk2.File(file.FullPath())

	if err != nil {
		return nil, err
	}

	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	sshUri := flagSet.String("ssh", "", "SSH")
	password := flagSet.String("pw", "", "Password")
	port := flagSet.Int("P", 22, "Port")
	sess := flagSet.String("load", "", "putty session")

	err = flagSet.Parse(strings.Split(lnk.StringData.CommandLineArguments, " "))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("parse \"%s\" flags", file.Name()))
	}

	ssh, err := url.Parse(fmt.Sprintf("ssh://%s", *sshUri))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("parsing %s", *sshUri))
	}

	name := file.Name()
	name = strings.TrimSuffix(name, ".lnk")

	sshAccess := accesses.NewSSH()
	sshAccess.SetName(name)
	sshAccess.SetPassword(*password)
	sshAccess.SetHost(ssh.Host)
	sshAccess.SetPort(*port)
	sshAccess.SetLogin(ssh.User.Username())

	if *sess != "" {
		sshAccess.SetSession(*sess)
	}

	return sshAccess, nil
}
