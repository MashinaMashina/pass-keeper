package puttyrun

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/pkg/iocustom"
)

func (lp *puttyRun) action(c *cli.Context) error {
	if lp.Config.String("putty.app") == "" {
		fmt.Fprintln(lp.Stdout, "Не указан адрес приложения Putty")
		fmt.Fprintln(lp.Stdout, "Введите адрес exe файла:")
		var path string
		err := iocustom.Fscanln(lp.Stdin, &path)

		if err != nil {
			return err
		}

		stat, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintln(lp.Stdout, "Нет такого файла")
			} else {
				fmt.Fprintln(lp.Stdout, "Неверный адрес приложения", err)
			}
			return nil
		}

		if stat.IsDir() {
			fmt.Fprintln(lp.Stdout, "Вы указали адрес папки, укажите адрес файла")
			return nil
		}

		lp.Config.Set("putty.app", path)
	}

	app := lp.Config.String("putty.app")

	var parameters []storage.Param
	parameters = append(parameters, params.NewEq("type", "ssh"))

	if c.Args().First() != "" {
		parameters = append(parameters, params.NewLike("name", c.Args().First()+"%"))
	}

	access, err := lp.Storage.FindOne(parameters...)
	if err != nil {
		return err
	}

	fmt.Fprintln(lp.Stdout, "Запуск "+access.Name())

	err = exec.Command(app, lp.argsFromAccess(access)...).Start()
	if err != nil {
		return err
	}

	return nil
}

func (lp *puttyRun) argsFromAccess(access accesstype.Access) (args []string) {
	if access.Host() != "" {
		host := access.Host()

		if access.Login() != "" {
			host = fmt.Sprintf("%s@%s", access.Login(), host)
		}

		args = append(args, "-ssh", host)
	}

	if access.Port() != 0 {
		args = append(args, "-P", fmt.Sprintf("%d", access.Port()))
	}

	if access.Password() != "" {
		args = append(args, "-pw", access.Password())
	}

	if access.Session() != "" {
		args = append(args, "--load", access.Session())
	}

	return
}
