package accessadd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/pkg/clibell"
	"pass-keeper/pkg/cliselect"
	"strconv"
	"time"
)

func (l *accessAdd) action(c *cli.Context) error {
	var value string
	var err error
	var access accesstype.Access

	if value = c.String("type"); value == "" {
		promt := cliselect.Select{
			Label:  "Choose type",
			Items:  []string{"ssh"},
			Stdout: clibell.Instance(l.Stdout),
			Mode:   cliselect.ModeSimple,
		}

		if l.DTO.Config.String("main.mode") == "interactive" {
			promt.Mode = cliselect.ModeInteractive
		}

		_, value, err = promt.Run()
		if err != nil {
			return err
		}

		if l.DTO.Config.String("main.mode") == "interactive" {
			// Без этого фикса promptui.Select съедает часть следующего вывода
			time.Sleep(time.Millisecond)
		}
	}

	switch value {
	case "ssh":
		access = accesstype.NewSSH()
	default:
		return errors.New(fmt.Sprintf("invalid type %s", value))
	}

	for {
		fmt.Fprint(l.Stdout, "Enter name: ")
		_, err = fmt.Fscanln(l.Stdin, &value)
		if err != nil {
			return err
		}

		rows, err := l.Storage.List(
			params.NewEq("name", value),
			params.NewEq("type", access.Type()),
		)
		if err != nil {
			return err
		}

		if len(rows) == 0 {
			access.SetName(value)
			break
		} else {
			fmt.Fprintln(l.Stdout, "Name is exists, choose another name")
		}
	}

	fmt.Fprint(l.Stdout, "Enter hostname: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetHost(value)
	}

	fmt.Fprint(l.Stdout, "Enter port: ")
	if access.Port() != 0 {
		fmt.Fprint(l.Stdout, " (by default ", access.Port(), ")")
	}
	fmt.Fprint(l.Stdout, ": ")

	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && (access.Port() == 0 || err.Error() != "unexpected newline") {
		return err
	}
	if value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		access.SetPort(port)
	}

	fmt.Fprint(l.Stdout, "Enter login: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetLogin(value)
	}

	fmt.Fprint(l.Stdout, "Enter password: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetPassword(value)
	}

	err = l.Storage.Add(access)
	if err != nil {
		return err
	}

	fmt.Fprintln(l.Stdout, "Success add")

	return nil
}
