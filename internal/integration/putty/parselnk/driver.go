package parselnk

import (
	"flag"
	"fmt"
	lnk2 "github.com/parsiya/golnk"
	"github.com/pkg/errors"
	"net/url"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/pkg/filesystem"
	"pass-keeper/pkg/flagcustom"
	"strings"
)

var ErrEmptyLink = errors.New("link has no data")

func (lp *linkParser) sshAccessByLnkFile(file filesystem.File) (accesstype.Access, error) {
	lnk, err := lnk2.File(file.FullPath())

	if err != nil {
		return nil, err
	}

	name := strings.TrimSuffix(file.Name(), ".lnk")
	name = lp.cleanFilename(name)

	return lp.accessByArguments(lnk.StringData.CommandLineArguments, name)
}

func (lp *linkParser) accessByArguments(args, name string) (accesstype.Access, error) {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	sshUri := flagSet.String("ssh", "", "SSH")
	password := flagSet.String("pw", "", "Password")
	port := flagSet.Int("P", 22, "Port")
	sess := flagSet.String("load", "", "putty session")

	flags, err := flagcustom.ParseFlags(args)
	if err != nil {
		return nil, err
	}

	err = flagSet.Parse(flags)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("parse \"%s\" flags", name))
	}

	ssh, err := url.Parse(fmt.Sprintf("ssh://%s", *sshUri))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("parsing %s", *sshUri))
	}

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
