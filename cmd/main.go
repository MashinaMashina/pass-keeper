package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty"
	"pass-keeper/internal/master"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{},
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	s, err := storage.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err = s.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	p := putty.New(s, cfg)
	m := master.New(s, cfg)

	app.Commands = append(app.Commands, p.Commands()...)
	app.Commands = append(app.Commands, m.Commands()...)

	err = cfg.LoadUserConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
