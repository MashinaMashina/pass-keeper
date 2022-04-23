package accesses

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/action/add"
	"pass-keeper/internal/accesses/action/edit"
	"pass-keeper/internal/accesses/action/list"
	"pass-keeper/internal/accesses/action/remove"
	"pass-keeper/internal/accesses/action/show"
	accessvalidate "pass-keeper/internal/accesses/action/validate"
	"pass-keeper/internal/app"
)

type access struct {
	app.DTO
}

func New(dto app.DTO) *access {
	return &access{dto}
}

func (p *access) Commands() []*cli.Command {
	commands := make([]*cli.Command, 0, 5)
	commands = append(commands, accesslist.New(p.DTO).Commands()...)
	commands = append(commands, accessadd.New(p.DTO).Commands()...)
	commands = append(commands, accessedit.New(p.DTO).Commands()...)
	commands = append(commands, accessshow.New(p.DTO).Commands()...)
	commands = append(commands, accessremove.New(p.DTO).Commands()...)
	commands = append(commands, accessvalidate.New(p.DTO).Commands()...)

	return []*cli.Command{
		{
			Name:        "access",
			Usage:       "Accesses",
			Subcommands: commands,
		},
	}
}
