package accessadd

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessAdd struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *accessAdd {
	a := &accessAdd{
		storage: s,
		config:  cfg,
	}

	return a
}

func (l *accessAdd) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:   "add",
		Usage:  "Добавить доступ",
		Action: l.action,
	})

	return commands
}
