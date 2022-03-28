package accessremove

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"strings"
)

func (l *accessRemove) action(c *cli.Context) error {
	var parameters []storage.Param

	if c.Args().First() != "" {
		parameters = append(parameters, params.NewLike("name", c.Args().First()+"%"))
	}

	row, err := l.storage.FindOne(parameters...)
	if err != nil {
		return err
	}

	var confirm string
	fmt.Print(fmt.Sprintf("Удалить %s? (Y/n) ", row.Name()))
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		return nil
	}

	err = l.storage.Remove(row)
	if err != nil {
		return err
	}

	fmt.Println("Удалено")

	return nil
}
