package test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"pass-keeper/internal"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/app"
	"pass-keeper/internal/config"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var stdin = os.Stdin
var stdout = os.Stdout

func buildTestApp() (app.DTO, *cli.App, error) {
	dto, err := testingDTO()
	if err != nil {
		return dto, nil, err
	}

	app := &cli.App{
		Commands: []*cli.Command{},
	}

	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true

	commands, err := internal.CollectCommands(dto)
	if err != nil {
		return dto, nil, errors.Wrap(err, "init commands")
	}

	app.Commands = append(app.Commands, commands...)

	return dto, app, nil
}

func testingDTO() (app.DTO, error) {
	dto := app.DTO{
		Stdout: stdout,
		Stdin:  stdin,
	}

	cfg := config.NewConfig()
	err := cfg.InitFromData([]byte("{\"master.password\":\"c4ca4238a0b923820dcc509a6f75849b\",\"storage.source\":\":memory:\"}"))
	if err != nil {
		return dto, errors.Wrap(err, "init config")
	}

	internal.FillConfig(cfg)

	dto.Config = cfg

	s, err := sqlite.New(cfg)
	if err != nil {
		return dto, errors.Wrap(err, "init storage")
	}

	dto.Storage = s

	return dto, nil
}

var boundary = []byte("==--delimiter--==")

type appTestCase struct {
	Args    []string
	Comment string
	Before  func(io.Writer, io.Reader, app.DTO)
	Check   func(*testing.T, app.DTO, string)
}

var dateFormat string

func appTestOutput(t *testing.T, cases []appTestCase) {
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
			testCase.Before(inWriter, outReader, dto)
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
		expect.Session() != real.Session() ||
		expect.Valid() != real.Valid() {
		t.Error(fmt.Sprintf("accesses not equal\nexpected: %+v,\nbut real: %+v", expect, real))
		return false
	}

	if !reflect.DeepEqual(expect.Params().All(), real.Params().All()) {
		t.Error(fmt.Sprintf("access params not equal\nexpected: %+v,\nbut real: %+v", expect.Params().All(), real.Params().All()))

		return false
	}

	return true
}

func FolderTest() string {
	_, testdata, _, _ := runtime.Caller(0)
	testdata = path.Dir(testdata) + "/testdata"

	return testdata
}
