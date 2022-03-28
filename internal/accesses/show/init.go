package accessshow

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessShow struct {
	storage      storage.Storage
	accessConfig *config.Part
}

func New(s storage.Storage, p *config.Part) *accessShow {
	a := &accessShow{
		storage:      s,
		accessConfig: p,
	}

	a.fillConfig()

	return a
}

func (l *accessShow) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:   "show",
		Usage:  "Информация о доступе",
		Action: l.action,
	})

	return commands
}
