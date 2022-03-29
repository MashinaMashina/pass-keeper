package basedriver

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage"
	"pass-keeper/internal/accesses/storage/params"
	"pass-keeper/internal/config"
	"pass-keeper/pkg/encrypt"
)

type BaseDriver struct {
	Db            *sql.DB
	Config        *config.Config
	StorageConfig *config.Part
	Key           []byte
}

func (s *BaseDriver) Add(access accesstype.Access) error {
	if access.Name() == "" {
		return fmt.Errorf("name of access can not be empty")
	}

	stmt, err := s.Db.Prepare("INSERT INTO accesses" +
		"(type, name, host, port, login, password, session, valid)" +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?)")

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

	_, err = stmt.Exec(access.Type(), access.Name(), access.Host(), access.Port(),
		login, password, access.Session(), access.Valid())

	if err != nil {
		return err
	}

	return nil
}

func (s *BaseDriver) Update(id int, access accesstype.Access) error {
	if access.Name() == "" {
		return fmt.Errorf("name of access can not be empty")
	}

	stmt, err := s.Db.Prepare("UPDATE accesses SET " +
		"type=?, name=?, host=?, port=?, login=?, password=?, session=?, valid=?" +
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

	_, err = stmt.Exec(access.Type(), access.Name(), access.Host(), access.Port(),
		login, password, access.Session(), access.Valid(), id)

	if err != nil {
		return err
	}

	return nil
}

func (s *BaseDriver) Save(access accesstype.Access) error {
	id, err := s.FindId(access)

	if err != nil {
		return err
	}

	if id > 0 {
		return s.Update(id, access)
	}

	return s.Add(access)
}

func (s *BaseDriver) Remove(access accesstype.Access) error {
	id, err := s.FindId(access)

	if err != nil {
		return err
	}

	if id > 0 {
		prepare, err := s.Db.Prepare("DELETE FROM accesses WHERE id=?")
		if err != nil {
			return err
		}
		defer prepare.Close()

		_, err = prepare.Exec(id)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("not found rows")
}

func (s *BaseDriver) Exists(access accesstype.Access) (bool, error) {
	id, err := s.FindId(access)

	if err != nil {
		return false, err
	}

	return id > 0, nil
}

func (s *BaseDriver) FindId(access accesstype.Access) (int, error) {
	stmt, err := s.Db.Prepare("SELECT id FROM accesses WHERE type=? AND name=? AND host=? LIMIT 1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(access.Type(), access.Name(), access.Host()).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return id, nil
}

func (s *BaseDriver) List(params ...storage.Param) ([]accesstype.Access, error) {
	query := squirrel.
		Select("type", "name", "host", "login", "port", "password").
		From("accesses")

	for _, param := range params {
		switch param.ParamType() {
		case "like":
			query = query.Where(param.Value()[0]+" LIKE ?", param.Value()[1])
		case "eq":
			query = query.Where(param.Value()[0]+" = ?", param.Value()[1])
		}
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
	var rows []accesstype.Access

	for stmt.Next() {
		access, err = s.decodeRow(stmt)
		if err != nil {
			return nil, err
		}

		rows = append(rows, access)
	}

	return rows, nil
}

func (s *BaseDriver) FindOne(parameters ...storage.Param) (accesstype.Access, error) {
	rows, err := s.List(parameters...)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("not found rows")
	}

	if len(rows) > 1 {
		names := make([]string, 0, len(rows))

		for _, row := range rows {
			names = append(names, row.Name())
		}

		promt := promptui.Select{
			Label: "Доступно несколько вариантов",
			Items: names,
		}

		_, res, err := promt.Run()
		if err != nil {
			return nil, err
		}

		return s.FindOne(append(parameters, params.NewEq("name", res))...)
	}

	return rows[0], nil
}

func (s *BaseDriver) decodeRow(rows *sql.Rows) (accesstype.Access, error) {
	var typo string
	var name string
	var host string
	var login string
	var port int
	var password string
	var access accesstype.Access
	var err error

	if err = rows.Scan(&typo, &name, &host, &login, &port, &password); err != nil {
		return nil, err
	}

	switch typo {
	case "ssh":
		access = accesstype.NewSSH()
	default:
		return nil, fmt.Errorf("unknown access type: " + typo)
	}

	login, err = s.decode(login)
	if err != nil {
		return nil, err
	}

	password, err = s.decode(password)
	if err != nil {
		return nil, err
	}

	access.SetHost(host)
	access.SetName(name)
	access.SetLogin(login)
	access.SetPort(port)
	access.SetPassword(password)

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
		masterPass, err := hex.DecodeString(s.Config.Part("master").Get("password"))

		if err != nil {
			panic("decoding hex of master password error: " + err.Error())
		}

		appKey, err := hex.DecodeString(s.Config.Part("main").Get("key"))
		if err != nil {
			return nil
		}

		s.Key = append(masterPass, appKey...)
	}

	return s.Key
}

func (s *BaseDriver) encode(data string) (string, error) {
	// Пробелами дополняем строку до 16 символов
	return encrypt.EncryptAES(s.getKey(), data)
}

func (s *BaseDriver) decode(data string) (string, error) {
	return encrypt.DecryptAES(s.getKey(), data)
}