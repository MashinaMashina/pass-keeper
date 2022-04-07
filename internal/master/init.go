package master

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type master struct {
	app.DTO
}

func New(dto app.DTO) *master {
	m := &master{dto}

	virtual := m.Config.Virtual()
	m.Config.SetVirtual(append(virtual, "master.password"))

	err := m.initConfig()
	if err != nil {
		panic(err)
	}

	return m
}

func (m *master) Commands() []*cli.Command {
	return []*cli.Command{}
}
