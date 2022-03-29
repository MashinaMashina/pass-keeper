package master

import (
	"pass-keeper/internal"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/driver/sqlite"
	"pass-keeper/internal/config"
)

func configAndStorage() (*config.Config, storage.Storage, error) {
	cfg, err := internal.NewConfig("")
	if err != nil {
		return nil, nil, err
	}

	storage, err := sqlite.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, storage, nil
}
