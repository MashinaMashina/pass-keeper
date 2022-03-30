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
	"strconv"
	"time"
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
		"(type, name, host, port, login, password, session, valid, created_at, updated_at)" +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

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

	res, err := stmt.Exec(access.Type(), access.Name(), access.Host(), access.Port(),
		login, password, access.Session(), access.Valid(), now.Unix(), now.Unix())

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
		"type=?, name=?, host=?, port=?, login=?, password=?, session=?, valid=?, updated_at=?" +
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

	_, err = stmt.Exec(access.Type(), access.Name(), access.Host(), access.Port(),
		login, password, access.Session(), access.Valid(), now.Unix(), access.ID())

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
		Select("id", "type", "name", "host", "login", "port", "password", "created_at", "updated_at").
		From("accesses")

	for _, param := range params {
		switch param.ParamType() {
		case "like":
			query = query.Where(param.Value()[0]+" LIKE ?", param.Value()[1])
		case "eq":
			query = query.Where(param.Value()[0]+" = ?", param.Value()[1])
		case "limit":
			i, _ := strconv.Atoi(param.Value()[0])
			query = query.Limit(uint64(i))
		default:
			return nil, fmt.Errorf("invalid param type %s", param.ParamType())
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
	var id int
	var typo string
	var name string
	var host string
	var login string
	var port int
	var password string
	var access accesstype.Access
	var err error
	var createdAt int64
	var updatedAt int64

	if err = rows.Scan(&id, &typo, &name, &host, &login, &port, &password, &createdAt, &updatedAt); err != nil {
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

	access.SetID(id)
	access.SetHost(host)
	access.SetName(name)
	access.SetLogin(login)
	access.SetPort(port)
	access.SetPassword(password)
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
