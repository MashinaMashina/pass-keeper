package accessvalidate

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/accesses/validate"
)

func (l *accessValidate) action(c *cli.Context) error {
	var parameters []storage.Param

	if c.Args().First() != "" {
		parameters = append(parameters, params.NewLike("name", c.Args().First()+"%"))
	}

	var rows []accesstype.Access

	if !c.Bool("all") {
		row, err := l.Storage.FindOne(parameters...)
		if err != nil {
			return err
		}

		rows = append(rows, row)
	} else {
		var err error
		rows, err = l.Storage.List(parameters...)

		if err != nil {
			return err
		}
	}

	verifyUnknownHosts := !c.Bool("allow-all-hosts")
	for _, row := range rows {
		l.validateRow(row, verifyUnknownHosts)
	}

	return nil
}

func (l *accessValidate) validateRow(access accesstype.Access, verifyUnknownHosts bool) error {
	fmt.Fprint(l.Stdout, "Checking ", access.Name(), ": ")

	var (
		valid bool
		err   error
	)

	switch access.Type() {
	case "ssh":
		valid, err = validate.ValidateSSH(access, verifyUnknownHosts, l.Stdout, l.Stdin)
	default:
		return errors.New("not found validator")
	}

	if valid {
		fmt.Fprintln(l.Stdout, "is valid")
	} else {
		fmt.Fprintln(l.Stdout, "not valid")
	}

	if err != nil {
		fmt.Fprintln(l.Stdout, "Error:", err)
	}

	if access.Valid() != valid {
		access.SetValid(valid)
		l.Storage.Save(access)
	}

	return nil
}
