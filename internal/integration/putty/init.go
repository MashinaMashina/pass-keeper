package putty

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty/parselnk"
)

type putty struct {
	storage     storage.Storage
	config      *config.Config
	puttyConfig *config.Part
}

func New(s storage.Storage, cfg *config.Config) *putty {
	h := config.NewPart()

	err := cfg.AddPart("putty", h)
	if err != nil {
		fmt.Println("error:", err.Error())
		return nil
	}

	return &putty{
		storage:     s,
		config:      cfg,
		puttyConfig: h,
	}
}

func (p *putty) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "putty",
			Usage:       "Putty интеграция",
			Subcommands: parselnk.New(p.storage, p.puttyConfig).Commands(),
		},
	}
}
