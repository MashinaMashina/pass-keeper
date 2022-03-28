package accessedit

import (
	"fmt"
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

	access, err := l.storage.FindOne(parameters...)
	if err != nil {
		return err
	}

	id, err := l.storage.FindId(access)
	if err != nil {
		return err
	}

	fmt.Println("Редактирование " + access.Name())
	fmt.Println("Если не хотите менять строку, пропускайте нажатием Enter")

	var value string
	var edit []string
	for {
		fmt.Print("Введите имя: ")
		_, err = fmt.Scanln(&value)
		if err != nil && err.Error() != "unexpected newline" {
			return err
		}

		if value != "" {
			rows, err := l.storage.List(
				params.NewEq("name", value),
				params.NewEq("type", access.Type()),
			)
			if err != nil {
				return err
			}

			if len(rows) == 0 {
				access.SetName(value)
				edit = append(edit, "имя")
				break
			} else {
				fmt.Println("Такое имя уже есть, выберите другое")
			}
		} else {
			break
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
		edit = append(edit, "хост")
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
		edit = append(edit, "порт")
	}

	fmt.Print("Введите логин: ")
	value = ""
	_, err = fmt.Scanln(&value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetLogin(value)
		edit = append(edit, "логин")
	}

	fmt.Print("Введите пароль: ")
	value = ""
	_, err = fmt.Scanln(&value)
	if err != nil && err.Error() != "unexpected newline" {
		return err
	}
	if value != "" {
		access.SetPassword(value)
		edit = append(edit, "пароль")
	}

	if len(edit) == 0 {
		fmt.Println("Нечего менять")
		return nil
	}

	fmt.Println(fmt.Sprintf("%+v", access))

	err = l.storage.Update(id, access)
	if err != nil {
		return err
	}

	fmt.Println("Обновлены поля: " + strings.Join(edit, ", "))

	return nil
}
