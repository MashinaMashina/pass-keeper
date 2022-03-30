package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"pass-keeper/internal"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
)

func main() {
	err := RunApp()
	if err != nil {
		log.Fatalln(err)
	}
}

func RunApp() error {
	cfg, err := internal.NewConfig("~/.pass-keeper.json")
	if err != nil {
		return err
	}

	storage, err := sqlite.New(cfg)
	if err != nil {
		return err
	}
	defer storage.Close()

	commands, err := internal.CollectCommands(storage, cfg)
	if err != nil {
		return err
	}

	err = cfg.LoadUserConfig()
	if err != nil {
		return err
	}

	err = cfg.Init()
	if err != nil {
		return err
	}

	app, err := AppBuild(commands)
	if err != nil {
		return err
	}

	err = app.Run(os.Args)
	if err != nil {
		return err
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
