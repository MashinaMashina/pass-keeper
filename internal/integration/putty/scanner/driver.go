package scanner

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

func (ls *linkScanner) sshAccessByLnkFile(file filesystem.File) (accesstype.Access, error) {
	lnk, err := lnk2.File(file.FullPath())

	if err != nil {
		return nil, err
	}

	name := strings.TrimSuffix(file.Name(), ".lnk")
	name = ls.cleanFilename(name)

	return ls.accessByArguments(lnk.StringData.CommandLineArguments, name, filesystem.PrettyPath(file.Dir()))
}

func (ls *linkScanner) accessByArguments(args, name, group string) (accesstype.Access, error) {
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
	sshAccess.SetGroup(group)

	if *sess != "" {
		sshAccess.Params().Set("session_name", *sess)
	}

	if sshAccess.Host() == "" && sshAccess.Params().Value("session_name") == "" {
		return nil, ErrEmptyLink
	}

	return sshAccess, nil
}
