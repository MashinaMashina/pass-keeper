package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func (s *sqlite) initConfig() error {
	source := s.BaseDriver.Config.String("storage.source")

	if source == "" {
		source = "~/.pass-keeper.db"
		s.BaseDriver.Config.Set("storage.source", source)
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
