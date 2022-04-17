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

	row, err := l.Storage.FindOne(parameters...)
	if err != nil {
		return err
	}

	var confirm string
	fmt.Fprint(l.Stdout, fmt.Sprintf("Удалить %s? (Y/n) ", row.Name()))
	fmt.Fscanln(l.Stdin, &confirm)

	if strings.ToLower(confirm) != "y" {
		return nil
	}

	err = l.Storage.Remove(row)
	if err != nil {
		return err
	}

	fmt.Fprintln(l.Stdout, "Deleted")

	return nil
}
