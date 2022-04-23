package accesslist

import (
	"fmt"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
)

func (l *accessList) action(c *cli.Context) error {
	parameters := make([]storage.Param, 0, 1)

	if c.Args().First() != "" {
		parameters = append(parameters, params.NewLike("name", c.Args().First()+"%"))
	}

	rows, err := l.Storage.List(parameters...)
	if err != nil {
		return err
	}

	if c.Bool("list") {
		return l.list(c, rows)
	} else {
		return l.table(c, rows)
	}
}

func (l *accessList) table(c *cli.Context, rows []accesstype.Access) error {
	if len(rows) == 0 {
		return nil
	}

	tbl := table.New("", "", "", "", "")

	tbl.WithWriter(l.Stdout)
	tbl.WithPadding(1)

	var line [3]string
	i := 0
	for _, row := range rows {
		line[i] = row.Name()
		i++

		if i == 3 {
			i = 0
			tbl.AddRow(line[0], line[1], line[2])
			line = [3]string{}
		}
	}

	if i != 0 {
		tbl.AddRow(line[0], line[1], line[2])
	}

	tbl.Print()

	return nil
}

func (l *accessList) list(c *cli.Context, rows []accesstype.Access) error {
	tbl := table.New("Type", "Name", "Host", "Login", "Updated")

	tbl.WithWriter(l.Stdout)
	tbl.WithPadding(1)

	for _, row := range rows {
		host := row.Host()
		if row.Port() != 0 {
			host = fmt.Sprintf("%s:%d", host, row.Port())
		}

		tbl.AddRow(row.Type(), row.Name(), host, row.Login(), row.UpdatedAt().Format(l.Config.String("main.date_format")))
	}

	tbl.Print()

	return nil
}
