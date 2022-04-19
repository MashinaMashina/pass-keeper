package accessshow

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type accessShow struct {
	app.DTO
}

func New(dto app.DTO) *accessShow {
	a := &accessShow{dto}

	return a
}

func (l *accessShow) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:   "show",
		Usage:  "Access info",
		Action: l.action,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "mask",
				Usage: "Search by mask. Example: %site.ru%",
			},
		},
	})

	return commands
}
