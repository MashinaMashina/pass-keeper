package internal

import (
	"encoding/json"
	"github.com/urfave/cli/v2"
	"pass-keeper/internal/accesses"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/config"
	"pass-keeper/internal/integration/putty"
	"pass-keeper/internal/master"
)

func NewConfig(file string) (*config.Config, error) {
	cfg, err := config.New(file)
	if err != nil {
		return nil, err
	}

	mainConf := config.NewPart()
	mainConf.Set("key", key)
	cfg.AddPart("main", mainConf)

	return cfg, nil
}

func CollectCommands(storage storage.Storage, cfg *config.Config) ([]*cli.Command, error) {
	var commands []*cli.Command

	commands = append(commands, accesses.New(storage, cfg).Commands()...)
	commands = append(commands, putty.New(storage, cfg).Commands()...)
	commands = append(commands, master.New(storage, cfg).Commands()...)

	return commands, nil
}

func TestingConfigAndStorage() (*config.Config, storage.Storage, error) {
	cfg, err := NewConfig("")
	if err != nil {
		return nil, nil, err
	}

	s, err := sqlite.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	_, err = CollectCommands(s, cfg)
	if err != nil {
		return nil, nil, err
	}

	jsonBytes := []byte("{\"master\":{\"file\":\"~/.pass-keeper.master\",\"password\":\"c4ca4238a0b923820dcc509a6f75849b\"},\"storage\":{\"source\":\":memory:\"}}")
	err = json.Unmarshal(jsonBytes, cfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, s, nil
}
