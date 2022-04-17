package parselnk

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type linkParser struct {
	app.DTO
}

func New(dto app.DTO) *linkParser {
	lp := &linkParser{dto}

	return lp
}

func (lp *linkParser) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:  "scan",
		Usage: "Collects accesses from putty shortcuts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Usage:       "Shortcut folder or shortcut file",
				Required:    false,
				Value:       "./",
				Destination: nil,
			},
		},
		Action: lp.action,
	})

	return commands
}
