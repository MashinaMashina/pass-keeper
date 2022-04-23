package runner

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type filezillaRun struct {
	app.DTO
}

func New(dto app.DTO) *filezillaRun {
	lp := &filezillaRun{dto}

	return lp
}

func (lp *filezillaRun) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run filezilla",
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
