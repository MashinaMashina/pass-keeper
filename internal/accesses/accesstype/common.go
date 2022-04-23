package accesstype

import (
	"time"
)

type Params interface {
	Fill(map[string]string)
	Value(string) string
	Exists(string) bool
	Set(string, string)
	All() map[string]string
}

type params struct {
	values map[string]string
}

func NewParams() Params {
	return &params{
		values: make(map[string]string),
	}
}

func (p *params) Fill(v map[string]string) {
	p.values = v
}

func (p *params) Value(k string) string {
	val, _ := p.values[k]

	return val
}

func (p *params) Exists(k string) bool {
	_, exists := p.values[k]

	return exists
}

func (p *params) Set(k string, v string) {
	p.values[k] = v
}

func (p *params) All() map[string]string {
	return p.values
}

type Access interface {
	ID() int
	SetID(id int)
	Type() string
	SetType(typo string)
	Name() string
	SetName(name string)
	Host() string
	SetHost(host string)
	Port() int
	SetPort(port int)
	Login() string
	SetLogin(login string)
	Password() string
	SetPassword(password string)
	Group() string
	SetGroup(group string)
	Valid() bool
	SetValid(valid bool)
	CreatedAt() time.Time
	SetCreatedAt(t time.Time)
	UpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
	Params() Params
}

type access struct {
	id        int
	typo      string
	name      string
	host      string
	port      int
	login     string
	password  string
	group     string
	valid     bool
	createdAt time.Time
	updatedAt time.Time
	params    Params
}

func NewUnknown() Access {
	return &access{params: NewParams()}
}

func (a *access) ID() int {
	return a.id
}

func (a *access) SetID(id int) {
	a.id = id
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

func (a *access) Group() string {
	return a.group
}

func (a *access) SetGroup(group string) {
	a.group = group
}

func (a *access) Valid() bool {
	return a.valid
}

func (a *access) SetValid(valid bool) {
	a.valid = valid
}

func (a *access) CreatedAt() time.Time {
	return a.createdAt
}

func (a *access) SetCreatedAt(t time.Time) {
	a.createdAt = t
}
func (a *access) UpdatedAt() time.Time {
	return a.createdAt
}

func (a *access) SetUpdatedAt(t time.Time) {
	a.createdAt = t
}

func (a *access) Params() Params {
	return a.params
}
