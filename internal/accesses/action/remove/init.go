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
	return []*cli.Command{
		{
			Name:  "remove",
			Usage: "Access remove",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "mask",
					Usage: "Search by mask. Example: %site.ru%",
				},
				&cli.BoolFlag{
					Name:    "all",
					Aliases: []string{"A"},
					Usage:   "validate all",
				},
			},
			Action: l.action,
		},
	}
}
