package internal

import (
	"pass-keeper/internal/config"
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
