package accesslist

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type accessList struct {
	app.DTO
}

func New(dto app.DTO) *accessList {
	a := &accessList{dto}

	return a
}

func (l *accessList) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:  "list",
		Usage: "Отображает список доступов",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
			},
		},
		Action: l.action,
	})

	return commands
}
