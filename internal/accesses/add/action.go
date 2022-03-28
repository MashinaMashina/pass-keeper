package accessadd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/params"
	"strconv"
	"time"
)

func (l *accessAdd) action(c *cli.Context) error {
	promt := promptui.Select{
		Label: "Выберите тип",
		Items: []string{"ssh"},
	}

	_, res, err := promt.Run()
	if err != nil {
		return err
	}

	// Без этого фикса promptui.Select съедает часть следующего вывода
	time.Sleep(time.Millisecond)

	var value string

	var access accesstype.Access
	switch res {
	case "ssh":
		access = accesstype.NewSSH()
	}

	for {
		fmt.Print("Введите имя: ")
		_, err = fmt.Scanln(&value)
		if err != nil {
			return err
		}

		rows, err := l.storage.List(
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
			fmt.Println("Такое имя уже есть, выберите другое")
		}
	}

	fmt.Print("Введите хост: ")
	value = ""
	_, err = fmt.Scanln(&value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetHost(value)
	}

	fmt.Print("Введите порт")
	if access.Port() != 0 {
		fmt.Print(" (по умолчанию ", access.Port(), ")")
	}
	fmt.Print(": ")

	value = ""
	_, err = fmt.Scanln(&value)
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

	fmt.Print("Введите логин: ")
	value = ""
	_, err = fmt.Scanln(&value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetLogin(value)
	}

	fmt.Print("Введите пароль: ")
	value = ""
	_, err = fmt.Scanln(&value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetPassword(value)
	}

	fmt.Println(fmt.Sprintf("%+v", access))

	err = l.storage.Add(access)
	if err != nil {
		return err
	}

	fmt.Println("Добавлено")

	return nil
}
