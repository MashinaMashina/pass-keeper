package parselnk

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type linkParser struct {
	storage       storage.Storage
	puttyConfig   *config.Part
	replaceConfig [][2]string
}

func New(s storage.Storage, p *config.Part) *linkParser {
	lp := &linkParser{
		storage:     s,
		puttyConfig: p,
	}

	lp.fillConfig()

	return lp
}

func (lp *linkParser) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:  "scan",
		Usage: "Собирает доступы с putty ярлыков",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Usage:       "Папка с ярлыками или файл ярлыка",
				Required:    false,
				Value:       "./",
				Destination: nil,
			},
		},
		Action: lp.action,
	})

	return commands
}
