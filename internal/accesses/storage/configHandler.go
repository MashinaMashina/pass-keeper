package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func (s *storage) fillHandler() {
	d := s.storageConfig.DefaultValues()
	d["source"] = "~/.pass-keeper.db"
	s.storageConfig.SetDefaultValues(d)

	i := s.storageConfig.InstallFields()
	i = append(i, "source")
	s.storageConfig.SetInstallFields(i)

	f := s.storageConfig.FieldNames()
	f["source"] = "файл для хранения доступов"
	s.storageConfig.SetFieldNames(f)

	s.storageConfig.SetValidate(s.validateConfig)
}

func (s *storage) validateConfig() error {
	source := s.storageConfig.Get("source")

	if source == "" {
		return fmt.Errorf("DB source not specified")
	}

	var err error
	source, err = homedir.Expand(source)

	if err != nil {
		return errors.Wrap(err, "Expand db source")
	}

	s.db, err = sql.Open("sqlite3", source)

	if err != nil {
		return errors.Wrap(err, "open db source")
	}

	return nil
}
