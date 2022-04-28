package txt

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
	"pass-keeper/internal/integration/txt/scanner"
)

type txt struct {
	app.DTO
}

func New(dto app.DTO) *txt {
	return &txt{dto}
}

func (p *txt) Commands() []*cli.Command {
	commands := make([]*cli.Command, 0, 1)

	commands = append(commands, scanner.New(p.DTO).Commands()...)

	return []*cli.Command{
		{
			Name:        "txt",
			Usage:       "text documents integration",
			Subcommands: commands,
			Action: func(c *cli.Context) error {
				return commands[0].Action(c)
			},
		},
	}
}
