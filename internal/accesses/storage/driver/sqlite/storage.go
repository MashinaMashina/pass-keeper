package sqlite

import (
	"io"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/driver/basedriver"
	"pass-keeper/internal/config"
)

type sqlite struct {
	basedriver.BaseDriver
}

func New(cfg *config.Config, stdin io.ReadCloser, stdout io.WriteCloser) (storage.Storage, error) {
	s := &sqlite{
		basedriver.BaseDriver{
			Config: cfg,
			Stdin:  stdin,
			Stdout: stdout,
		},
	}

	err := s.initConfig()
	if err != nil {
		return nil, err
	}

	return s, nil
}
