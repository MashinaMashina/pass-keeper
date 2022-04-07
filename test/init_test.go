package test

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"pass-keeper/internal"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/app"
	"pass-keeper/internal/config"
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

	s, err := sqlite.New(cfg, dto.Stdin, dto.Stdout)
	if err != nil {
		return dto, errors.Wrap(err, "init storage")
	}

	dto.Storage = s

	return dto, nil
}
