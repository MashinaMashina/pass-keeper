package runner

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type puttyRun struct {
	app.DTO
}

func New(dto app.DTO) *puttyRun {
	lp := &puttyRun{dto}

	return lp
}

func (lp *puttyRun) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run putty",
			Action: lp.action,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "mask",
					Usage: "Search by mask. Example: %site.ru%",
				},
			},
		},
	}
}
