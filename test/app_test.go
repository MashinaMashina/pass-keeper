package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/app"
	"strings"
	"testing"
)

func TestAppCRUD(t *testing.T) {
	var cases []appTestCase

	access := accesstype.NewSSH()
	access.SetName("some1")
	access.SetHost("host.com")
	access.SetPort(33)
	access.SetLogin("login")
	access.SetPassword("qwerty123")

	cases = append(cases, appTestCase{
		Args: []string{"access", "list"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			assert.Equal(t, "", strings.TrimSpace(output), "")
		},
	})

	cases = append(cases, appTestCase{
		Args: []string{"access", "list", "-l"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			assert.Equalf(t, "Type Name Host Login Updated", strings.TrimSpace(output), "")
		},
	})

	cases = append(cases, appTestCase{
		Args: []string{"access", "add", "--type", "ssh"},
		Before: func(w io.Writer, r io.Reader, dto app.DTO) {
			fmt.Fprintln(w, access.Name())
			fmt.Fprintln(w, access.Host())
			fmt.Fprintln(w, access.Port())
			fmt.Fprintln(w, access.Login())
			fmt.Fprintln(w, access.Password())
		},
		Check: func(t *testing.T, dto app.DTO, output string) {
			if assert.Equalf(t, "Введите имя: Введите хост: Введите порт (по умолчанию 22): "+
				"Введите логин: Введите пароль: Добавлено", strings.TrimSpace(output), "") {
				rows, err := dto.Storage.List(params.NewEq("name", access.Name()))
				if err != nil {
					t.Error(err)
					return
				}

				if len(rows) == 1 {
					access.SetCreatedAt(rows[0].CreatedAt())
					access.SetUpdatedAt(rows[0].UpdatedAt())
				}

				equalOneAccess(t, rows, access)
			}
		},
	})

	cases = append(cases, appTestCase{
		Args: []string{"access", "show", access.Name()},
		Check: func(t *testing.T, dto app.DTO, output string) {
			output = cleanLines(output)

			expect := fmt.Sprintf("Имя %s\n"+
				"Хост %s\n"+
				"Порт %d\n"+
				"Логин %s\n"+
				"Пароль %s\n"+
				"Имя сессии\n"+
				"Валиден Нет\n"+
				"Добавлен %s\n"+
				"Изменен %s",
				access.Name(), access.Host(), access.Port(), access.Login(), access.Password(),
				access.CreatedAt().Format(dateFormat),
				access.UpdatedAt().Format(dateFormat),
			)

			assert.Equalf(t, expect, strings.TrimSpace(cleanLines(output)), "")
		},
	})

	cases = append(cases, appTestCase{
		Args: []string{"access", "list"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			assert.Equalf(t, access.Name(), strings.TrimSpace(output), "")
		},
	})

	cases = append(cases, appTestCase{
		Args: []string{"access", "list", "-l"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			expect := fmt.Sprintf("Type Name Host Login Updated\n%s %s %s:%d %s %s",
				access.Type(), access.Name(), access.Host(), access.Port(),
				access.Login(), access.UpdatedAt().Format(dateFormat))

			assert.Equalf(t, expect, strings.TrimSpace(cleanLines(output)), "")
		},
	})

	cases = append(cases, appTestCase{
		Args:    []string{"access", "edit", access.Name()},
		Comment: "без изменений",
		Before: func(w io.Writer, r io.Reader, dto app.DTO) {
			fmt.Fprint(w, "\n\n\n\n\n")
		},
		Check: func(t *testing.T, dto app.DTO, output string) {
			except := fmt.Sprintf("Редактирование %s\nЕсли не хотите менять строку, пропускайте нажатием Enter\n"+
				"Введите имя: Введите хост: Введите порт: Введите логин: Введите пароль: Нечего менять",
				access.Name())

			if assert.Equalf(t, except, strings.TrimSpace(output), "") {
				rows, err := dto.Storage.List(params.NewEq("name", access.Name()))
				if err != nil {
					t.Error(err)
					return
				}

				equalOneAccess(t, rows, access)
			}
		},
	})

	access2 := accesstype.NewSSH()
	access2.SetName("some2")
	access2.SetHost("host2.com")
	access2.SetPort(44)
	access2.SetLogin("login1")
	access2.SetPassword("qwerty123!")

	cases = append(cases, appTestCase{
		Args:    []string{"access", "edit", access.Name()},
		Comment: "с изменениями",
		Before: func(w io.Writer, r io.Reader, dto app.DTO) {
			fmt.Fprintln(w, access2.Name())
			fmt.Fprintln(w, access2.Host())
			fmt.Fprintln(w, access2.Port())
			fmt.Fprintln(w, access2.Login())
			fmt.Fprintln(w, access2.Password())
		},
		Check: func(t *testing.T, dto app.DTO, output string) {
			expect := fmt.Sprintf("Редактирование %s\n"+
				"Если не хотите менять строку, пропускайте нажатием Enter\n"+
				"Введите имя: Введите хост: Введите порт: Введите логин: Введите пароль: Обновлены поля: имя, хост, порт, логин, пароль", access.Name())

			if assert.Equalf(t, expect, strings.TrimSpace(output), "") {
				rows, err := dto.Storage.List(params.NewEq("name", access2.Name()))
				if err != nil {
					t.Error(err)
					return
				}

				if len(rows) == 1 {
					access2.SetCreatedAt(rows[0].CreatedAt())
					access2.SetUpdatedAt(rows[0].UpdatedAt())
				}

				equalOneAccess(t, rows, access2)
			}
		},
	})

	cases = append(cases, appTestCase{
		Args:    []string{"access", "remove", access2.Name()},
		Comment: "не подтверждаем удаление",
		Before: func(w io.Writer, r io.Reader, dto app.DTO) {
			fmt.Fprintln(w, "n")
		},
		Check: func(t *testing.T, dto app.DTO, output string) {
			expect := fmt.Sprintf("Удалить %s? (Y/n)", access2.Name())

			if assert.Equalf(t, expect, strings.TrimSpace(output), "") {
				rows, err := dto.Storage.List(params.NewEq("name", access2.Name()))
				if err != nil {
					t.Error(err)
					return
				}

				equalOneAccess(t, rows, access2)
			}
		},
	})

	cases = append(cases, appTestCase{
		Args:    []string{"access", "remove", access2.Name()},
		Comment: "подтверждаем удаление",
		Before: func(w io.Writer, r io.Reader, dto app.DTO) {
			fmt.Fprintln(w, "y")
		},
		Check: func(t *testing.T, dto app.DTO, output string) {
			expect := fmt.Sprintf("Удалить %s? (Y/n) Удалено", access2.Name())

			if assert.Equalf(t, expect, strings.TrimSpace(output), "") {
				rows, err := dto.Storage.List(params.NewEq("name", access2.Name()))
				if err != nil {
					t.Error(err)
					return
				}

				if len(rows) != 0 {
					t.Error("access not removed")
					return
				}
			}
		},
	})

	appTestOutput(t, cases)
}
