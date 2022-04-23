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
	return []*cli.Command{
		{
			Name:  "list",
			Usage: "Access list",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Value:   false,
				},
			},
			Action: l.action,
		},
	}
}
