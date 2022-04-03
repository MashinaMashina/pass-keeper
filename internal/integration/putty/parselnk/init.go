package parselnk

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type linkParser struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *linkParser {
	lp := &linkParser{
		storage: s,
		config:  cfg,
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
