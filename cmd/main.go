package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"pass-keeper/internal"
	"pass-keeper/internal/accesses"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty"
	"pass-keeper/internal/master"
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

	commands, err := AppPreBuild(storage, cfg)
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

func AppPreBuild(storage storage.Storage, cfg *config.Config) ([]*cli.Command, error) {
	var commands []*cli.Command

	commands = append(commands, accesses.New(storage, cfg).Commands()...)
	commands = append(commands, putty.New(storage, cfg).Commands()...)
	commands = append(commands, master.New(storage, cfg).Commands()...)

	return commands, nil
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
