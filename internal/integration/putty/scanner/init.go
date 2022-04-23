package scanner

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type linkScanner struct {
	app.DTO
}

func New(dto app.DTO) *linkScanner {
	lp := &linkScanner{dto}

	return lp
}

func (ls *linkScanner) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "scan",
			Usage:  "Collects accesses from putty shortcuts",
			Action: ls.action,
		},
	}
}
