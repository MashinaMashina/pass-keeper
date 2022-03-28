package accessedit

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessEdit struct {
	storage      storage.Storage
	accessConfig *config.Part
}

func New(s storage.Storage, p *config.Part) *accessEdit {
	a := &accessEdit{
		storage:      s,
		accessConfig: p,
	}

	a.fillConfig()

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
