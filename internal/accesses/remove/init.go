package accessremove

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type accessRemove struct {
	storage      storage.Storage
	accessConfig *config.Part
}

func New(s storage.Storage, p *config.Part) *accessRemove {
	a := &accessRemove{
		storage:      s,
		accessConfig: p,
	}

	a.fillConfig()

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
