package accessremove

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessRemove struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *accessRemove {
	a := &accessRemove{
		storage: s,
		config:  cfg,
	}

	return a
}

func (l *accessRemove) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:   "remove",
		Usage:  "Удаление доступа",
		Action: l.action,
	})

	return commands
}
