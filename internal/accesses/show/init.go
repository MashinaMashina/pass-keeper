package accessshow

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessShow struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *accessShow {
	a := &accessShow{
		storage: s,
		config:  cfg,
	}

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
