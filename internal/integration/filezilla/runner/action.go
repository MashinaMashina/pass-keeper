package runner

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"net/url"
	"os"
	"os/exec"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	app2 "pass-keeper/internal/app"
	"pass-keeper/pkg/iocustom"
)

func (lp *filezillaRun) action(c *cli.Context) error {
	if lp.Config.String("filezilla.app") == "" {
		fmt.Fprintln(lp.Stdout, "FileZilla application address not specified")
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

		lp.Config.Set("filezilla.app", path)
	}

	app := lp.Config.String("filezilla.app")

	var parameters []storage.Param

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

	if access.Type() != "ssh" && access.Type() != "ftp" {
		return errors.New(fmt.Sprintf("filezilla not supported %s type", access.Type()))
	}

	fmt.Fprintln(lp.Stdout, "Starting "+access.Name())

	err = exec.Command(app, lp.argsFromAccess(access)...).Start()
	if err != nil {
		return err
	}

	return nil
}

func (lp *filezillaRun) argsFromAccess(access accesstype.Access) (args []string) {
	var u string
	switch access.Type() {
	case "ssh":
		u = "sftp://"
	case "ftp":
		u = "ftp://"
	}

	if access.Login() != "" {
		u += url.QueryEscape(access.Login())

		if access.Password() != "" {
			u += ":" + url.QueryEscape(access.Password())
		}

		u += "@"
	}

	u += fmt.Sprintf("%s:%d", url.QueryEscape(access.Host()), access.Port())

	args = append(args, u)

	return
}
