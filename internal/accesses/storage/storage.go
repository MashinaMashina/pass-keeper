package storage

import (
	"database/sql"
	"fmt"
	"os"
	"pass-keeper/internal/accesses"
	"pass-keeper/internal/config"
)

type storage struct {
	db            *sql.DB
	storageConfig *config.Part
}

func New(cfg *config.Config) (Storage, error) {
	part := config.NewPart()
	err := cfg.AddPart("storage", part)

	if err != nil {
		return nil, err
	}

	s := &storage{
		storageConfig: part,
	}

	s.fillHandler()

	return s, nil
}

func (s *storage) Add(access accesses.Access) error {

}

func (s *storage) Update(access accesses.Access) error {
	//TODO implement me
	panic("implement me")
}

func (s *storage) Remove(access accesses.Access) error {
	//TODO implement me
	panic("implement me")
}

func (s *storage) Exists(name string) (accesses.Access, error) {
	//TODO implement me
	panic("implement me")
}

func (s *storage) List() ([]accesses.Access, error) {
	fmt.Println(os.UserHomeDir())

	return []accesses.Access{}, nil
}

func (s *storage) Close() error {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
