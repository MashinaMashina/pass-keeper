package accessvalidate

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type accessValidate struct {
	app.DTO
}

func New(dto app.DTO) *accessValidate {
	a := &accessValidate{dto}

	return a
}

func (l *accessValidate) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "validate",
			Usage: "Check access valid",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "allow-all-hosts",
					Usage: "allow all unknown hosts",
				},
				&cli.BoolFlag{
					Name:    "all",
					Aliases: []string{"A"},
					Usage:   "validate all",
				},
				&cli.BoolFlag{
					Name:  "mask",
					Usage: "Search by mask. Example: %site.ru%",
				},
			},
			Action: l.action,
		},
	}
}
