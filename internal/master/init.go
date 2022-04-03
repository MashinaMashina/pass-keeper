package master

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/config"
)

type master struct {
	storage storage.Storage
	config  *config.Config
}

func New(s storage.Storage, cfg *config.Config) *master {
	m := &master{
		storage: s,
		config:  cfg,
	}

	virtual := m.config.Virtual()
	m.config.SetVirtual(append(virtual, "master.password"))

	err := m.initConfig()
	if err != nil {
		panic(err)
	}

	return m
}

func (m *master) Commands() []*cli.Command {
	return []*cli.Command{}
}
