package storage

import (
	"github.com/Masterminds/squirrel"
	"pass-keeper/internal/accesses/accesstype"
)

type Storage interface {
	Add(access accesstype.Access) error
	Update(access accesstype.Access) error
	Save(access accesstype.Access) error
	Remove(access accesstype.Access) error
	Exists(access accesstype.Access) (bool, error)
	FindExists(access accesstype.Access) (accesstype.Access, error)
	List(...Param) ([]accesstype.Access, error)
	Close() error
}

type Param func(*squirrel.SelectBuilder)
