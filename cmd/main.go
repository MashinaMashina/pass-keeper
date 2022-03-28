package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"pass-keeper/internal/accesses"
	"pass-keeper/internal/accesses/storage/sqlite"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty"
	"pass-keeper/internal/master"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{},
	}

	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true

	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err = storage.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	app.Commands = append(app.Commands, accesses.New(storage, cfg).Commands()...)
	app.Commands = append(app.Commands, putty.New(storage, cfg).Commands()...)
	app.Commands = append(app.Commands, master.New(storage, cfg).Commands()...)

	err = cfg.LoadUserConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = cfg.Init()
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
