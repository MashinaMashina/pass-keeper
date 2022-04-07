package accesses

import (
	"github.com/urfave/cli/v2"
	accessadd "pass-keeper/internal/accesses/add"
	accessedit "pass-keeper/internal/accesses/edit"
	accesslist "pass-keeper/internal/accesses/list"
	accessremove "pass-keeper/internal/accesses/remove"
	accessshow "pass-keeper/internal/accesses/show"
	"pass-keeper/internal/app"
)

type access struct {
	app.DTO
}

func New(dto app.DTO) *access {
	return &access{dto}
}

func (p *access) Commands() []*cli.Command {
	commands := accesslist.New(p.DTO).Commands()
	commands = append(commands, accessadd.New(p.DTO).Commands()...)
	commands = append(commands, accessedit.New(p.DTO).Commands()...)
	commands = append(commands, accessshow.New(p.DTO).Commands()...)
	commands = append(commands, accessremove.New(p.DTO).Commands()...)

	return []*cli.Command{
		{
			Name:        "access",
			Usage:       "Доступы",
			Subcommands: commands,
		},
	}
}
