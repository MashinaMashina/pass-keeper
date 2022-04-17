package accessedit

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"strconv"
	"strings"
)

func (l *accessEdit) action(c *cli.Context) error {
	var parameters []storage.Param

	if c.Args().First() != "" {
		parameters = append(parameters, params.NewLike("name", c.Args().First()+"%"))
	}

	access, err := l.Storage.FindOne(parameters...)
	if err != nil {
		return errors.Wrap(err, "get access from storage")
	}

	fmt.Fprintln(l.Stdout, "Editing "+access.Name())
	fmt.Fprintln(l.Stdout, "If you do not want to change the line, skip by pressing Enter")

	var value string
	var edit []string
	for {
		fmt.Fprint(l.Stdout, "Enter name: ")
		_, err = fmt.Fscanln(l.Stdin, &value)
		if err != nil && err.Error() != "unexpected newline" {
			return errors.Wrap(err, "scan name")
		}

		if value != "" {
			rows, err := l.Storage.List(
				params.NewEq("name", value),
				params.NewEq("type", access.Type()),
			)
			if err != nil {
				return errors.Wrap(err, "check access name is exists")
			}

			if len(rows) == 0 {
				access.SetName(value)
				edit = append(edit, "name")
				break
			} else {
				fmt.Fprintln(l.Stdout, "Name is exists, choose another name")
			}
		} else {
			break
		}
	}

	fmt.Fprint(l.Stdout, "Enter hostname: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return errors.Wrap(err, "scan host")
	}
	if value != "" {
		access.SetHost(value)
		edit = append(edit, "hostname")
	}

	fmt.Fprint(l.Stdout, "Enter port: ")

	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && (access.Port() == 0 || err.Error() != "unexpected newline") {
		return errors.Wrap(err, "scan port")
	}
	if value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		access.SetPort(port)
		edit = append(edit, "port")
	}

	fmt.Fprint(l.Stdout, "Enter login: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return errors.Wrap(err, "scan login")
	}
	if value != "" {
		access.SetLogin(value)
		edit = append(edit, "login")
	}

	fmt.Fprint(l.Stdout, "Enter password: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return errors.Wrap(err, "scan password")
	}
	if value != "" {
		access.SetPassword(value)
		edit = append(edit, "password")
	}

	if len(edit) == 0 {
		fmt.Fprintln(l.Stdout, "Nothing to change")
		return nil
	}

	err = l.Storage.Update(access)
	if err != nil {
		return errors.Wrap(err, "update access")
	}

	fmt.Fprintln(l.Stdout, "Updated fields: "+strings.Join(edit, ", "))

	return nil
}
