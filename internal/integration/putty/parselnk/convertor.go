package parselnk

import (
	"flag"
	"fmt"
	lnk2 "github.com/parsiya/golnk"
	"github.com/pkg/errors"
	"net/url"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/pkg/filesystem"
	"strings"
)

var ErrEmptyLink = errors.New("link has no data")

func (lp *linkParser) sshAccessByLnkFile(file filesystem.File) (accesstype.Access, error) {
	lnk, err := lnk2.File(file.FullPath())

	if err != nil {
		return nil, err
	}

	return lp.sshAccessByLnk(lnk, file.Name())
}

func (lp *linkParser) sshAccessByLnk(lnk lnk2.LnkFile, name string) (accesstype.Access, error) {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	sshUri := flagSet.String("ssh", "", "SSH")
	password := flagSet.String("pw", "", "Password")
	port := flagSet.Int("P", 22, "Port")
	sess := flagSet.String("load", "", "putty session")

	err := flagSet.Parse(strings.Split(lnk.StringData.CommandLineArguments, " "))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("parse \"%s\" flags", name))
	}

	ssh, err := url.Parse(fmt.Sprintf("ssh://%s", *sshUri))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("parsing %s", *sshUri))
	}

	name = strings.TrimSuffix(name, ".lnk")
	name = lp.cleanFilename(name)

	sshAccess := accesstype.NewSSH()
	sshAccess.SetName(name)
	sshAccess.SetPassword(*password)
	sshAccess.SetHost(ssh.Host)
	sshAccess.SetPort(*port)
	sshAccess.SetLogin(ssh.User.Username())

	if *sess != "" {
		sshAccess.SetSession(*sess)
	}

	if sshAccess.Host() == "" && sshAccess.Session() == "" {
		return nil, ErrEmptyLink
	}

	return sshAccess, nil
}
