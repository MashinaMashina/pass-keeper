package accessadd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/pkg/clibell"
	"strconv"
	"time"
)

func (l *accessAdd) action(c *cli.Context) error {
	var value string
	var err error
	var access accesstype.Access

	if value = c.String("type"); value == "" {
		promt := promptui.Select{
			Label:  "Выберите тип",
			Items:  []string{"ssh"},
			Stdout: clibell.Instance(l.Stdout),
		}

		_, value, err = promt.Run()
		if err != nil {
			return err
		}

		// Без этого фикса promptui.Select съедает часть следующего вывода
		time.Sleep(time.Millisecond)
	}

	switch value {
	case "ssh":
		access = accesstype.NewSSH()
	default:
		return errors.New(fmt.Sprintf("invalid type %s", value))
	}

	for {
		fmt.Fprint(l.Stdout, "Введите имя: ")
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
			fmt.Fprintln(l.Stdout, "Такое имя уже есть, выберите другое")
		}
	}

	fmt.Fprint(l.Stdout, "Введите хост: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetHost(value)
	}

	fmt.Fprint(l.Stdout, "Введите порт")
	if access.Port() != 0 {
		fmt.Fprint(l.Stdout, " (по умолчанию ", access.Port(), ")")
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

	fmt.Fprint(l.Stdout, "Введите логин: ")
	value = ""
	_, err = fmt.Fscanln(l.Stdin, &value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetLogin(value)
	}

	fmt.Fprint(l.Stdout, "Введите пароль: ")
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

	fmt.Fprintln(l.Stdout, "Добавлено")

	return nil
}
