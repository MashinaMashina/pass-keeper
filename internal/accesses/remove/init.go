package accessremove

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type accessRemove struct {
	app.DTO
}

func New(dto app.DTO) *accessRemove {
	a := &accessRemove{dto}

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
