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
	return []*cli.Command{
		{
			Name:   "add",
			Usage:  "Add access",
			Action: l.action,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "type",
					Usage:       "Access type",
					Required:    false,
					Destination: nil,
				},
			},
		},
	}
}
