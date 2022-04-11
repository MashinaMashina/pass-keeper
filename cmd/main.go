package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"pass-keeper/internal"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/app"
	"pass-keeper/internal/config"
)

func main() {
	err := RunApp()
	if err != nil {
		log.Fatalln(err)
	}
}

func RunApp() error {
	dto := app.DTO{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
	}

	cfg := config.NewConfig()
	err := cfg.InitFromFile(internal.ConfigFile())
	if err != nil {
		return errors.Wrap(err, "init config")
	}
	defer cfg.SaveToFile()

	dto.Config = cfg

	internal.FillConfig(dto.Config)

	storage, err := sqlite.New(dto.Config, dto.Stdin, dto.Stdout)
	if err != nil {
		return errors.Wrap(err, "init storage")
	}
	defer storage.Close()

	dto.Storage = storage

	commands, err := internal.CollectCommands(dto)
	if err != nil {
		return errors.Wrap(err, "init commands")
	}

	application := &cli.App{
		Commands: []*cli.Command{},
		Reader:   dto.Stdin,
		Writer:   dto.Stdout,
	}

	application.EnableBashCompletion = true
	application.UseShortOptionHandling = true

	application.Commands = append(application.Commands, commands...)

	err = application.Run(os.Args)
	if err != nil {
		return errors.Wrap(err, "app runtime")
	}

	return nil
}
