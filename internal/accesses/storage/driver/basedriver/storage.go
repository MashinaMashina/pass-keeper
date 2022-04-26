package basedriver

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/config"
	"pass-keeper/pkg/encrypt"
	"time"
)

type BaseDriver struct {
	Db     *sql.DB
	Config *config.Config
	Key    []byte
}

func (s *BaseDriver) Add(access accesstype.Access) error {
	if access.Name() == "" {
		return fmt.Errorf("name of access can not be empty")
	}

	stmt, err := s.Db.Prepare("INSERT INTO accesses" +
		"(`type`, `name`, `host`, `port`, `login`, `password`, `group`, `valid`, `created_at`, `updated_at`, `params`)" +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return errors.Wrap(err, "prepare add access")
	}

	login, err := s.encode(access.Login())
	if err != nil {
		return errors.Wrap(err, "encoding login")
	}
	password, err := s.encode(access.Password())
	if err != nil {
		return errors.Wrap(err, "encoding password")
	}

	now := time.Now()

	params, err := json.Marshal(access.Params().All())
	if err != nil {
		return err
	}

	res, err := stmt.Exec(access.Type(), access.Name(), access.Host(), access.Port(),
		login, password, access.Group(), access.Valid(), now.Unix(), now.Unix(), string(params))

	if err != nil {
		return err
	}

	access.SetCreatedAt(now)
	access.SetUpdatedAt(now)

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	access.SetID(int(id))

	return nil
}

func (s *BaseDriver) Update(access accesstype.Access) error {
	if access.Name() == "" {
		return fmt.Errorf("name of access can not be empty")
	}
	if access.ID() == 0 {
		return fmt.Errorf("ID of access can not be empty")
	}

	stmt, err := s.Db.Prepare("UPDATE accesses SET " +
		"`type`=?, `name`=?, `host`=?, `port`=?, `login`=?, `password`=?," +
		"`group`=?, `valid`=?, `updated_at`=?, `params`=?" +
		"WHERE id=?")

	if err != nil {
		return errors.Wrap(err, "prepare update access")
	}

	login, err := s.encode(access.Login())
	if err != nil {
		return errors.Wrap(err, "encoding login")
	}
	password, err := s.encode(access.Password())
	if err != nil {
		return errors.Wrap(err, "encoding password")
	}

	now := time.Now()

	params, err := json.Marshal(access.Params().All())
	if err != nil {
		return err
	}

	_, err = stmt.Exec(access.Type(), access.Name(), access.Host(), access.Port(),
		login, password, access.Group(), access.Valid(), now.Unix(), string(params), access.ID())

	if err != nil {
		return err
	}

	access.SetUpdatedAt(now)

	return nil
}

func (s *BaseDriver) Save(access accesstype.Access) error {
	if access.ID() > 0 {
		return s.Update(access)
	}

	return s.Add(access)
}

func (s *BaseDriver) Remove(access accesstype.Access) error {
	if access.ID() == 0 {
		return fmt.Errorf("ID of access can not be empty")
	}

	prepare, err := s.Db.Prepare("DELETE FROM accesses WHERE id=?")
	if err != nil {
		return err
	}
	defer prepare.Close()

	_, err = prepare.Exec(access.ID())
	if err != nil {
		return err
	}

	access.SetID(0)

	return nil
}

func (s *BaseDriver) Exists(access accesstype.Access) (bool, error) {
	return access.ID() > 0, nil
}

func (s *BaseDriver) FindExists(access accesstype.Access) (accesstype.Access, error) {
	row, err := s.List(params.NewEq("name", access.Name()), params.NewEq("type", access.Type()),
		params.NewEq("host", access.Host()), params.NewLimit(1))

	if err != nil {
		return nil, err
	}

	if len(row) == 0 {
		return nil, fmt.Errorf("not find rows")
	}

	return row[0], nil
}

func (s *BaseDriver) List(params ...storage.Param) ([]accesstype.Access, error) {
	query := squirrel.
		Select("id", "type", "name", "host", "login", "port", "password",
			"`group`", "valid", "created_at", "updated_at", "params").
		From("accesses")

	for _, param := range params {
		param(&query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	stmt, err := s.Db.Query(sql, args...)

	if err != nil {
		return nil, errors.Wrap(err, "get accesses from DB")
	}
	defer stmt.Close()

	var access accesstype.Access
	rows := make([]accesstype.Access, 0, 16)

	for stmt.Next() {
		access, err = s.decodeRow(stmt)
		if err != nil {
			return nil, err
		}

		rows = append(rows, access)
	}

	return rows, nil
}

func (s *BaseDriver) decodeRow(rows *sql.Rows) (accesstype.Access, error) {
	var (
		id         int
		typo       string
		name       string
		host       string
		login      string
		port       int
		password   string
		group      string
		valid      bool
		access     accesstype.Access
		err        error
		createdAt  int64
		updatedAt  int64
		parameters string
	)

	if err = rows.Scan(&id, &typo, &name, &host, &login, &port, &password,
		&group, &valid, &createdAt, &updatedAt, &parameters); err != nil {
		return nil, errors.Wrap(err, "scanning storage data to variables")
	}

	if val, exists := accesstype.Types[typo]; exists {
		access = val()
	} else {
		return nil, errors.New(fmt.Sprintf("invalid type %s", typo))
	}

	login, err = s.decode(login)
	if err != nil {
		return nil, errors.Wrap(err, "login decoding")
	}

	password, err = s.decode(password)
	if err != nil {
		return nil, errors.Wrap(err, "password decoding")
	}

	var p map[string]string
	err = json.Unmarshal([]byte(parameters), &p)
	if err != nil {
		return nil, err
	}

	access.Params().Fill(p)

	access.SetID(id)
	access.SetHost(host)
	access.SetName(name)
	access.SetLogin(login)
	access.SetPort(port)
	access.SetPassword(password)
	access.SetGroup(group)
	access.SetValid(valid)
	access.SetCreatedAt(time.Unix(createdAt, 0))
	access.SetUpdatedAt(time.Unix(updatedAt, 0))

	return access, nil
}

func (s *BaseDriver) Close() error {
	if s.Db != nil {
		err := s.Db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BaseDriver) getKey() []byte {
	if s.Key == nil {
		masterPass, err := hex.DecodeString(s.Config.String("master.password"))

		if err != nil {
			panic("decoding hex of master password error: " + err.Error())
		}

		appKey, err := hex.DecodeString(s.Config.String("main.key"))
		if err != nil {
			panic("invalid app key")
		}

		s.Key = append(masterPass, appKey...)
	}

	return s.Key
}

func (s *BaseDriver) encode(data string) (string, error) {
	res, err := encrypt.EncryptAES(s.getKey(), data)
	if err != nil {
		err = errors.Wrap(err, "encode storage data")
	}

	return res, err
}

func (s *BaseDriver) decode(data string) (string, error) {
	res, err := encrypt.DecryptAES(s.getKey(), data)
	if err != nil {
		err = errors.Wrap(err, "decode storage data (may be invalid master password?)")
	}

	return res, err
}
