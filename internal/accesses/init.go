package accesses

import (
	"github.com/urfave/cli/v2"
	accessadd "pass-keeper/internal/accesses/add"
	accessedit "pass-keeper/internal/accesses/edit"
	accesslist "pass-keeper/internal/accesses/list"
	accessremove "pass-keeper/internal/accesses/remove"
	accessshow "pass-keeper/internal/accesses/show"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type access struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *access {
	return &access{
		storage: s,
		config:  cfg,
	}
}

func (p *access) Commands() []*cli.Command {
	commands := accesslist.New(p.storage, p.config).Commands()
	commands = append(commands, accessadd.New(p.storage, p.config).Commands()...)
	commands = append(commands, accessedit.New(p.storage, p.config).Commands()...)
	commands = append(commands, accessshow.New(p.storage, p.config).Commands()...)
	commands = append(commands, accessremove.New(p.storage, p.config).Commands()...)

	return []*cli.Command{
		{
			Name:        "access",
			Usage:       "Доступы",
			Subcommands: commands,
		},
	}
}
