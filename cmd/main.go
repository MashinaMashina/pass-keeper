package main

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"pass-keeper/internal"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/config"
)

func main() {
	err := RunApp()
	if err != nil {
		log.Fatalln(err)
	}

}

func RunApp() error {
	cfg := config.NewConfig()
	err := cfg.InitFromFile("~/.pass-keeper.json")
	if err != nil {
		return errors.Wrap(err, "init config")
	}
	defer cfg.SaveToFile()

	internal.FillKey(cfg)

	storage, err := sqlite.New(cfg)
	if err != nil {
		return errors.Wrap(err, "init storage")
	}
	defer storage.Close()

	commands, err := internal.CollectCommands(storage, cfg)
	if err != nil {
		return errors.Wrap(err, "init commands")
	}

	app, err := AppBuild(commands)
	if err != nil {
		return errors.Wrap(err, "build app")
	}

	err = app.Run(os.Args)
	if err != nil {
		return errors.Wrap(err, "app runtime")
	}

	return nil
}

func AppBuild(commands []*cli.Command) (*cli.App, error) {
	app := &cli.App{
		Commands: []*cli.Command{},
	}

	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true

	app.Commands = append(app.Commands, commands...)

	return app, nil
}
