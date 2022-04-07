package test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/app"
	"regexp"
	"strings"
	"testing"
)

var boundary = []byte("==--delimiter--==")

type appCrudTestCase struct {
	Args    []string
	Comment string
	Before  func(io.Writer, io.Reader)
	Check   func(*testing.T, app.DTO, string)
}

func TestAppCRUD(t *testing.T) {
	var cases []appCrudTestCase

	access := accesstype.NewSSH()
	access.SetName("some1")
	access.SetHost("host.com")
	access.SetPort(33)
	access.SetLogin("login")
	access.SetPassword("qwerty123")

	cases = append(cases, appCrudTestCase{
		Args: []string{"access", "list"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			assert.Equal(t, "", strings.TrimSpace(output), "")
		},
	})

	cases = append(cases, appCrudTestCase{
		Args: []string{"access", "list", "-l"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			assert.Equalf(t, "Type Name Host Login Updated", strings.TrimSpace(output), "")
		},
	})

	cases = append(cases, appCrudTestCase{
		Args: []string{"access", "add", "--type", "ssh"},
		Before: func(w io.Writer, r io.Reader) {
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

	cases = append(cases, appCrudTestCase{
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

	cases = append(cases, appCrudTestCase{
		Args: []string{"access", "list"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			assert.Equalf(t, access.Name(), strings.TrimSpace(output), "")
		},
	})

	cases = append(cases, appCrudTestCase{
		Args: []string{"access", "list", "-l"},
		Check: func(t *testing.T, dto app.DTO, output string) {
			expect := fmt.Sprintf("Type Name Host Login Updated\n%s %s %s:%d %s %s",
				access.Type(), access.Name(), access.Host(), access.Port(),
				access.Login(), access.UpdatedAt().Format(dateFormat))

			assert.Equalf(t, expect, strings.TrimSpace(cleanLines(output)), "")
		},
	})

	cases = append(cases, appCrudTestCase{
		Args:    []string{"access", "edit", access.Name()},
		Comment: "без изменений",
		Before: func(w io.Writer, r io.Reader) {
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

	cases = append(cases, appCrudTestCase{
		Args:    []string{"access", "edit", access.Name()},
		Comment: "с изменениями",
		Before: func(w io.Writer, r io.Reader) {
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

	cases = append(cases, appCrudTestCase{
		Args:    []string{"access", "remove", access2.Name()},
		Comment: "не подтверждаем удаление",
		Before: func(w io.Writer, r io.Reader) {
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

	cases = append(cases, appCrudTestCase{
		Args:    []string{"access", "remove", access2.Name()},
		Comment: "подтверждаем удаление",
		Before: func(w io.Writer, r io.Reader) {
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

	appCRUDTestOutput(t, cases)
}

var dateFormat string

func appCRUDTestOutput(t *testing.T, cases []appCrudTestCase) {
	in, inWriter, err := os.Pipe()
	if err != nil {
		return
	}
	defer inWriter.Close()

	outReader, out, err := os.Pipe()
	if err != nil {
		return
	}

	bufscan := bufio.NewScanner(outReader)
	bufscan.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.Index(data, boundary); i >= 0 {
			return i + len(boundary), data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})

	stdout = out
	stdin = in

	dto, app, err := buildTestApp()
	if err != nil {
		t.Error(err)
		return
	}

	defer dto.Storage.Close()

	dateFormat = dto.Config.String("main.date_format")

	args := []string{os.Args[0]}
	for _, testCase := range cases {
		msg := "Run "
		msg += "\"" + strings.Join(testCase.Args, " ") + "\""

		if testCase.Comment != "" {
			msg += fmt.Sprintf(" (%s)", testCase.Comment)
		}

		t.Log(msg)

		if testCase.Before != nil {
			testCase.Before(inWriter, outReader)
		}

		err = app.Run(append(args, testCase.Args...))
		if err != nil {
			t.Error(errors.Wrap(err, "app result"))
			return
		}

		out.Write(boundary)

		if testCase.Check != nil && bufscan.Scan() {
			testCase.Check(t, dto, bufscan.Text())
		}

	}

	out.Close()
}

var regexpSpaces = regexp.MustCompile(`[\t\f\r ]{2,}|\t`)

func cleanLines(s string) string {
	s = regexpSpaces.ReplaceAllString(s, " ")
	lines := strings.Split(s, "\n")

	res := ""
	for _, line := range lines {
		res += strings.TrimSpace(line) + "\n"
	}

	return res
}

func equalOneAccess(t *testing.T, expect []accesstype.Access, real accesstype.Access) bool {
	if len(expect) == 0 {
		t.Error("not found rows")
		return false
	}

	if len(expect) > 1 {
		t.Error(fmt.Sprintf("too many rows (%d > 1)", len(expect)))
		return false
	}

	return equalAccess(t, expect[0], real)
}

func equalAccess(t *testing.T, expect, real accesstype.Access) bool {
	if expect.Name() != real.Name() ||
		expect.Host() != real.Host() ||
		expect.Port() != real.Port() ||
		expect.Login() != real.Login() ||
		expect.Password() != real.Password() ||
		expect.Session() != real.Session() {
		t.Error(fmt.Sprintf("accesses not equal\nexpected: %+v,\nbut real %+v", expect, real))

		return false
	}

	return true
}
