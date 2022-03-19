package storage

import "pass-keeper/internal/accesses"

type Storage interface {
	Add(access accesses.Access) error
	Update(access accesses.Access) error
	Remove(access accesses.Access) error
	Exists(name string) (accesses.Access, error)
	List() ([]accesses.Access, error)
	Close() error
}
