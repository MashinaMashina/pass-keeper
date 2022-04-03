package accessedit

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessEdit struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *accessEdit {
	a := &accessEdit{
		storage: s,
		config:  cfg,
	}

	return a
}

func (l *accessEdit) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:   "edit",
		Usage:  "Изменить доступ",
		Action: l.action,
	})

	return commands
}
