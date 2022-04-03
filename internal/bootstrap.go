package internal

import (
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty"
	"pass-keeper/internal/master"
)

func FillKey(cfg *config.Config) {
	virtual := cfg.Virtual()
	cfg.SetVirtual(append(virtual, "main.key"))

	cfg.Set("main.key", key)
}

func CollectCommands(storage storage.Storage, cfg *config.Config) ([]*cli.Command, error) {
	var commands []*cli.Command

	commands = append(commands, accesses.New(storage, cfg).Commands()...)
	commands = append(commands, putty.New(storage, cfg).Commands()...)
	commands = append(commands, master.New(storage, cfg).Commands()...)

	return commands, nil
}

func TestingConfigAndStorage() (*config.Config, storage.Storage, error) {
	cfg := config.NewConfig()
	err := cfg.InitFromData([]byte("{\"master.file\":\"~/.pass-keeper.master\",\"storage.source\":\":memory:\"}"))
	if err != nil {
		return nil, nil, err
	}

	FillKey(cfg)

	s, err := sqlite.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	_, err = CollectCommands(s, cfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, s, nil
}
