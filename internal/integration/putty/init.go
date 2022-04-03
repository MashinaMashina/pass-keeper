package putty

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty/parselnk"
	"pass-keeper/internal/integration/putty/puttyrun"
)

type putty struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *putty {
	return &putty{
		storage: s,
		config:  cfg,
	}
}

func (p *putty) Commands() []*cli.Command {
	var commands []*cli.Command
	commands = append(commands, parselnk.New(p.storage, p.config).Commands()...)
	commands = append(commands, puttyrun.New(p.storage, p.config).Commands()...)

	return []*cli.Command{
		{
			Name:        "putty",
			Usage:       "Putty интеграция",
			Subcommands: commands,
		},
	}
}
