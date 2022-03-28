package accesslist

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessList struct {
	storage      storage.Storage
	accessConfig *config.Part
}

func New(s storage.Storage, p *config.Part) *accessList {
	a := &accessList{
		storage:      s,
		accessConfig: p,
	}

	a.fillConfig()

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
