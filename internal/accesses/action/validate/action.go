package accessvalidate

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/accesses/validate"
	"pass-keeper/internal/app"
)

func (l *accessValidate) action(c *cli.Context) error {
	parameters := make([]storage.Param, 0, 1)

	if c.Args().First() != "" {
		if c.Bool("mask") {
			parameters = append(parameters, params.NewLike("name", c.Args().First()))
		} else {
			parameters = append(parameters, params.NewLike("name", c.Args().First()+"%"))
		}
	}

	var rows []accesstype.Access

	if !c.Bool("all") {
		row, err := app.FindOne(l.DTO, parameters...)
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

func (l *accessValidate) validateRow(access accesstype.Access, verifyUnknownHosts bool) {
	fmt.Fprint(l.Stdout, "Checking ", access.Name(), ": ")

	var (
		valid bool
		err   error
	)

	if val, exists := validate.Validators[access.Type()]; exists {
		valid, err = val(access, verifyUnknownHosts, l.Stdout, l.Stdin)
	} else {
		fmt.Fprintln(l.Stdout, "not found validator for type", access.Type())
		return
	}

	if valid {
		fmt.Fprintln(l.Stdout, "is valid")
	} else {
		fmt.Fprintln(l.Stdout, "not valid")
	}

	if err != nil {
		fmt.Fprintln(l.Stdout, "Error:", err)
		return
	}

	if access.Valid() != valid {
		access.SetValid(valid)
		l.Storage.Save(access)
	}

	return
}
