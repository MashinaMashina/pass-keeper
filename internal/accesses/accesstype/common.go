package accesstype

import (
	"fmt"
)

type access struct {
	typo      string
	name      string
	host      string
	port      int
	login     string
	password  string
	session   string
	valid     bool
	validated bool
}

type Access interface {
	Type() string
	SetType(typo string)
	Name() string
	SetName(typo string)
	Host() string
	SetHost(host string)
	Port() int
	SetPort(port int)
	Login() string
	SetLogin(login string)
	Password() string
	SetPassword(password string)
	Session() string
	SetSession(session string)
	Valid() bool
	Validate() error
}

func New() Access {
	return &access{typo: "unknown"}
}

func (a *access) Type() string {
	return a.typo
}

func (a *access) SetType(typo string) {
	a.typo = typo
}

func (a *access) Name() string {
	return a.name
}

func (a *access) SetName(name string) {
	a.name = name
}

func (a *access) Host() string {
	return a.host
}

func (a *access) SetHost(host string) {
	a.host = host
}

func (a *access) Port() int {
	return a.port
}

func (a *access) SetPort(port int) {
	a.port = port
}

func (a *access) Login() string {
	return a.login
}

func (a *access) SetLogin(login string) {
	a.login = login
}

func (a *access) Password() string {
	return a.password
}

func (a *access) SetPassword(password string) {
	a.password = password
}

func (a *access) Session() string {
	return a.session
}

func (a *access) SetSession(session string) {
	a.session = session
}

func (a *access) Valid() bool {
	return a.valid
}

func (a *access) Validate() error {
	return fmt.Errorf("validator not setted for accesstype %s", a.typo)
}
