package accessremove

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/app"
	"strings"
)

func (l *accessRemove) action(c *cli.Context) error {
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

	if len(rows) == 0 {
		return errors.New("not found rows")
	}

	var names string
	for _, row := range rows {
		names += row.Name() + ", "
	}
	names = names[:len(names)-2]

	var confirm string
	fmt.Fprint(l.Stdout, fmt.Sprintf("Remove %s? (Y/n) ", names))
	fmt.Fscanln(l.Stdin, &confirm)

	if strings.ToLower(confirm) != "y" {
		return nil
	}

	for _, row := range rows {
		err := l.Storage.Remove(row)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(l.Stdout, "Deleted")

	return nil
}
