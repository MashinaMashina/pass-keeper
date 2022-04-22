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
	app2 "pass-keeper/internal/app"
	"pass-keeper/pkg/iocustom"
)

func (lp *puttyRun) action(c *cli.Context) error {
	if lp.Config.String("putty.app") == "" {
		fmt.Fprintln(lp.Stdout, "Putty application address not specified")
		fmt.Fprintln(lp.Stdout, "Enter the address of the exe file:")
		var path string
		err := iocustom.Fscanln(lp.Stdin, &path)

		if err != nil {
			return err
		}

		stat, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintln(lp.Stdout, "No such file")
			} else {
				fmt.Fprintln(lp.Stdout, "Invalid application address", err)
			}
			return nil
		}

		if stat.IsDir() {
			fmt.Fprintln(lp.Stdout, "You have specified the folder address, specify the file address")
			return nil
		}

		lp.Config.Set("putty.app", path)
	}

	app := lp.Config.String("putty.app")

	var parameters []storage.Param
	parameters = append(parameters, params.NewEq("type", "ssh"))

	if c.Args().First() != "" {
		if c.Bool("mask") {
			parameters = append(parameters, params.NewLike("name", c.Args().First()))
		} else {
			parameters = append(parameters, params.NewLike("name", c.Args().First()+"%"))
		}
	}

	access, err := app2.FindOne(lp.DTO, parameters...)
	if err != nil {
		return err
	}

	fmt.Fprintln(lp.Stdout, "Starting "+access.Name())

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

	if access.Params().Exists("session_name") {
		args = append(args, "--load", access.Params().Value("session_name"))
	}

	return
}
