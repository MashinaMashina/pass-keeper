package parselnk

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type linkParser struct {
	storage storage.Storage
}

func New(s storage.Storage, h *config.Part) *linkParser {
	fillHandler(h)

	return &linkParser{
		storage: s,
	}
}

func (lp *linkParser) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:  "scan",
		Usage: "Собирает доступы с putty ярлыков",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "folder",
				Aliases:     []string{"f"},
				Usage:       "Папка с ярлыками",
				Required:    false,
				Value:       "./",
				Destination: nil,
			},
		},
		Action: lp.cliAction,
	})

	return commands
}
