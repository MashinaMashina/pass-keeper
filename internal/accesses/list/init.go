package accesslist

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessList struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *accessList {
	a := &accessList{
		storage: s,
		config:  cfg,
	}

	return a
}

func (l *accessList) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:  "list",
		Usage: "Отображает список доступов",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
			},
		},
		Action: l.action,
	})

	return commands
}
