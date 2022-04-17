package accessadd

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type accessAdd struct {
	app.DTO
}

func New(dto app.DTO) *accessAdd {
	a := &accessAdd{dto}

	return a
}

func (l *accessAdd) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:   "add",
		Usage:  "Добавить доступ",
		Action: l.action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "type",
				Usage:       "Тип доступа",
				Required:    false,
				Destination: nil,
			},
		},
	})

	return commands
}
