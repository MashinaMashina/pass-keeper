package accessshow

import (
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
)

func (l *accessShow) action(c *cli.Context) error {
	var likeParam storage.Param

	if c.Args().First() != "" {
		likeParam = params.NewLike("name", c.Args().First()+"%")
	}

	row, err := l.Storage.FindOne(likeParam)
	if err != nil {
		return err
	}

	return l.show(row)
}

func (l *accessShow) show(row accesstype.Access) error {
	tbl := table.New("", "")

	tbl.WithWriter(l.Stdout)
	tbl.WithPadding(1)

	isValid := "Нет"
	if row.Valid() {
		isValid = "Да"
	}

	tbl.AddRow("Name", row.Name())
	tbl.AddRow("Host", row.Host())
	tbl.AddRow("Port", row.Port())
	tbl.AddRow("Login", row.Login())
	tbl.AddRow("Password", row.Password())
	tbl.AddRow("Session name", row.Session())
	tbl.AddRow("Valid", isValid)
	tbl.AddRow("Added", row.CreatedAt().Format(l.Config.String("main.date_format")))
	tbl.AddRow("Updated", row.UpdatedAt().Format(l.Config.String("main.date_format")))

	tbl.Print()

	return nil
}
