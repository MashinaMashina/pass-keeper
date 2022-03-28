package storage

import (
	"pass-keeper/internal/accesses/accesstype"
)

type Storage interface {
	Add(access accesstype.Access) error
	Update(id int, access accesstype.Access) error
	Save(access accesstype.Access) error
	Remove(access accesstype.Access) error
	Exists(accesstype.Access) (bool, error)
	List(...Param) ([]accesstype.Access, error)
	FindId(access accesstype.Access) (int, error)
	FindOne(...Param) (accesstype.Access, error)
	Close() error
}

type Param interface {
	ParamType() string
	Value() []string
}
