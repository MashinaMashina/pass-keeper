package sqlite

import (
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/driver/basedriver"
	"pass-keeper/internal/config"
)

type sqlite struct {
	basedriver.BaseDriver
}

func New(cfg *config.Config) (storage.Storage, error) {
	part := config.NewPart()
	err := cfg.AddPart("storage", part)

	if err != nil {
		return nil, err
	}

	s := &sqlite{
		basedriver.BaseDriver{
			StorageConfig: part,
			Config:        cfg,
		},
	}

	s.fillConfig()

	return s, nil
}
