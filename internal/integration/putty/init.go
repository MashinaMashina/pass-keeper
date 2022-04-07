package putty

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
	"pass-keeper/internal/integration/putty/parselnk"
	"pass-keeper/internal/integration/putty/puttyrun"
)

type putty struct {
	app.DTO
}

func New(dto app.DTO) *putty {
	return &putty{dto}
}

func (p *putty) Commands() []*cli.Command {
	var commands []*cli.Command
	commands = append(commands, parselnk.New(p.DTO).Commands()...)
	commands = append(commands, puttyrun.New(p.DTO).Commands()...)

	return []*cli.Command{
		{
			Name:        "putty",
			Usage:       "Putty интеграция",
			Subcommands: commands,
		},
	}
}
