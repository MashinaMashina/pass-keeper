package puttyrun

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/app"
)

type puttyRun struct {
	app.DTO
}

func New(dto app.DTO) *puttyRun {
	lp := &puttyRun{dto}

	return lp
}

func (lp *puttyRun) Commands() []*cli.Command {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:   "run",
		Usage:  "Run putty",
		Action: lp.action,
	})

	return commands
}
