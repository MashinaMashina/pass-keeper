package internal

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses"
	"pass-keeper/internal/app"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty"
	"pass-keeper/internal/master"
)

func FillConfig(cfg *config.Config) {
	cfg.SetVirtual(append(cfg.Virtual(), "main.key"))

	cfg.Set("main.key", key)
	cfg.Set("main.date_format", "15:04 02/01/06")
}

func ConfigFile() string {
	return configFile
}

func CollectCommands(dto app.DTO) ([]*cli.Command, error) {
	var commands []*cli.Command

	commands = append(commands, accesses.New(dto).Commands()...)
	commands = append(commands, putty.New(dto).Commands()...)
	commands = append(commands, master.New(dto).Commands()...)

	return commands, nil
}
