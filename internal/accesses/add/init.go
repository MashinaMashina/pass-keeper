package accessadd

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessAdd struct {
	storage      storage.Storage
	accessConfig *config.Part
}

func New(s storage.Storage, p *config.Part) *accessAdd {
	a := &accessAdd{
		storage:      s,
		accessConfig: p,
	}

	a.fillConfig()

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
