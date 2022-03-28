package master

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type master struct {
	storage      storage.Storage
	config       *config.Config
	masterConfig *config.Part
}

func New(s storage.Storage, cfg *config.Config) *master {
	h := config.NewPart()

	err := cfg.AddPart("master", h)
	if err != nil {
		fmt.Println("error:", err.Error())
		return nil
	}

	m := &master{
		storage:      s,
		config:       cfg,
		masterConfig: h,
	}

	m.fillConfig()

	return m
}

func (m *master) Commands() []*cli.Command {
	return []*cli.Command{}
}
