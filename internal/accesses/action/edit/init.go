package accessedit

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type accessEdit struct {
	app.DTO
}

func New(dto app.DTO) *accessEdit {
	a := &accessEdit{dto}

	return a
}

func (l *accessEdit) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:  "edit",
		Usage: "Edit access",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "mask",
				Usage: "Search by mask. Example: %site.ru%",
			},
		},
		Action: l.action,
	})

	return commands
}
