package putty

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
	"pass-keeper/internal/integration/putty/runner"
	"pass-keeper/internal/integration/putty/scanner"
)

type putty struct {
	app.DTO
}

func New(dto app.DTO) *putty {
	return &putty{dto}
}

func (p *putty) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, runner.New(p.DTO).Commands()...)
	commands = append(commands, scanner.New(p.DTO).Commands()...)

	return []*cli.Command{
		{
			Name:        "putty",
			Usage:       "Putty integration",
			Subcommands: commands,
			Action: func(c *cli.Context) error {
				return commands[0].Action(c)
			},
		},
	}
}
