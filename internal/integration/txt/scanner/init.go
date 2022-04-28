package scanner

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/app"
)

type scanner struct {
	app.DTO
	access   accesstype.Access
	accesses []accesstype.Access
	filled   map[string]struct{}
	nextType string
	comment  string
	group    string
	name     string
}

func New(dto app.DTO) *scanner {
	lp := &scanner{
		DTO: dto,
	}

	return lp
}

func (s *scanner) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "scan",
			Usage:  "Collects accesses from putty shortcuts",
			Action: s.action,
		},
	}
}
