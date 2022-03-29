package sqlite

import (
	"github.com/pkg/errors"
)

func (s *sqlite) migration() error {
	_, err := s.BaseDriver.Db.Exec("CREATE TABLE IF NOT EXISTS `accesses` (" +
		"id INTEGER  PRIMARY KEY AUTOINCREMENT," +
		"type TEXT NULL," +
		"name TEXT NULL," +
		"host TEXT NULL," +
		"port INTEGER NULL," +
		"login TEXT NULL," +
		"password TEXT NULL," +
		"session TEXT NULL," +
		"valid INT DEFAULT 0" +
		")")
	if err != nil {
		return errors.Wrap(err, "Migration of tables")
	}

	return nil
}
