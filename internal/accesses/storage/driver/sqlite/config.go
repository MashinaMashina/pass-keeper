package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func (s *sqlite) fillConfig() {
	d := s.BaseDriver.StorageConfig.DefaultValues()
	d["source"] = "~/.pass-keeper.db"
	s.BaseDriver.StorageConfig.SetDefaultValues(d)

	i := s.BaseDriver.StorageConfig.InstallFields()
	i = append(i, "source")
	s.BaseDriver.StorageConfig.SetInstallFields(i)

	f := s.BaseDriver.StorageConfig.FieldNames()
	f["source"] = "файл для хранения доступов"
	s.BaseDriver.StorageConfig.SetFieldNames(f)

	s.BaseDriver.StorageConfig.SetInit(s.initConfig)
}

func (s *sqlite) initConfig() error {
	source := s.BaseDriver.StorageConfig.Get("source")

	if source == "" {
		return fmt.Errorf("DB source not specified")
	}

	var err error
	source, err = homedir.Expand(source)

	if err != nil {
		return errors.Wrap(err, "Expand db source")
	}

	s.BaseDriver.Db, err = sql.Open("sqlite3", source)

	if err != nil {
		return errors.Wrap(err, "open db source")
	}

	return s.migration()
}
