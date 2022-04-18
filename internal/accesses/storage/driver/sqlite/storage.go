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
	s := &sqlite{
		basedriver.BaseDriver{
			Config: cfg,
		},
	}

	err := s.initConfig()
	if err != nil {
		return nil, err
	}

	return s, nil
}
