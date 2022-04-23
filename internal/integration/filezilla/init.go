package filezilla

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
	"pass-keeper/internal/integration/filezilla/runner"
)

type filezilla struct {
	app.DTO
}

func New(dto app.DTO) *filezilla {
	return &filezilla{dto}
}

func (p *filezilla) Commands() []*cli.Command {
	var commands []*cli.Command
	commands = append(commands, runner.New(p.DTO).Commands()...)

	return []*cli.Command{
		{
			Name:        "filezilla",
			Usage:       "FileZilla integration",
			Subcommands: commands,
			Action: func(c *cli.Context) error {
				return commands[0].Action(c)
			},
		},
	}
}
