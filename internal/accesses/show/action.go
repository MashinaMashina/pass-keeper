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

	row, err := l.storage.FindOne(likeParam)
	if err != nil {
		return err
	}

	return l.show(row)
}

func (l *accessShow) show(row accesstype.Access) error {
	tbl := table.New("", "")

	tbl.WithPadding(1)

	isValid := "Нет"
	if row.Valid() {
		isValid = "Да"
	}

	tbl.AddRow("Имя", row.Name())
	tbl.AddRow("Хост", row.Host())
	tbl.AddRow("Порт", row.Port())
	tbl.AddRow("Логин", row.Login())
	tbl.AddRow("Пароль", row.Password())
	tbl.AddRow("Имя сессии", row.Session())
	tbl.AddRow("Валиден", isValid)

	tbl.Print()

	return nil
}
