package test

import (
	"pass-keeper/internal/accesses/accesstype"
	"pass-keeper/internal/accesses/storage/params"
	"testing"
)

func TestStorageCRUD(t *testing.T) {
	dto, err := testingDTO()
	defer dto.Storage.Close()
	if err != nil {
		t.Error(err)
		return
	}

	access := accesstype.NewSSH()
	access.SetName("testing")
	access.SetHost("host.com")
	access.SetLogin("user")
	access.SetPassword("qwerty123")

	if access.ID() != 0 {
		t.Error("invalid id of new record:", access.ID())
		return
	}

	t.Log("add access")
	err = dto.Storage.Add(access)
	if err != nil {
		t.Error(err)
		return
	}

	insertId := access.ID()

	t.Log("Insert ID:", insertId)

	if insertId != 1 {
		t.Error("invalid id of new record")
		return
	}

	t.Log("getting rows by parameters")
	rows, err := dto.Storage.List(params.NewEq("name", access.Name()))
	if err != nil {
		t.Error(err)
		return
	}

	if len(rows) != 1 {
		t.Error("wrong rows length 1!=", len(rows))
		return
	}

	if rows[0].ID() != insertId {
		t.Errorf("selected ID not equal inserted ID %d!=%d", rows[0].ID(), insertId)
		return
	}

	if rows[0].Name() != access.Name() ||
		rows[0].Host() != access.Host() ||
		rows[0].Login() != access.Login() ||
		rows[0].Password() != access.Password() ||
		rows[0].Port() != access.Port() {
		t.Error("invalid row in db")
		return
	}

	access.SetHost("new." + access.Host())

	t.Log("save access")
	err = dto.Storage.Save(access)
	if err != nil {
		t.Error(err)
		return
	}

	updatedId := access.ID()

	t.Log("updated ID:", updatedId)

	if updatedId != insertId {
		t.Errorf("invalid id of new record %d!=%d", updatedId, insertId)
		return
	}

	t.Log("getting rows by parameters")
	rows, err = dto.Storage.List(params.NewEq("name", access.Name()))
	if err != nil {
		t.Error(err)
		return
	}

	if len(rows) != 1 {
		t.Error("wrong rows length 1!=", len(rows))
		return
	}

	if rows[0].ID() != updatedId {
		t.Errorf("selected ID not equal updated ID %d!=%d", rows[0].ID(), insertId)
		return
	}

	if rows[0].Name() != access.Name() ||
		rows[0].Host() != access.Host() ||
		rows[0].Login() != access.Login() ||
		rows[0].Password() != access.Password() ||
		rows[0].Port() != access.Port() {
		t.Error("invalid row in db")
		return
	}

	t.Log("remove access")
	err = dto.Storage.Remove(access)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("getting all rows")
	rows, err = dto.Storage.List()
	if err != nil {
		t.Error(err)
		return
	}

	if len(rows) != 0 {
		t.Error("error with removing")
		return
	}
}

func TestMultipleRows(t *testing.T) {
	dto, err := testingDTO()
	defer dto.Storage.Close()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("add 3 rows")

	access1 := accesstype.NewSSH()
	access1.SetName("Какое-то имя 1")

	access2 := accesstype.NewSSH()
	access2.SetName("Какое-то имя 2")

	access3 := accesstype.NewSSH()
	access3.SetName("Какое-то имя 3")

	err = dto.Storage.Add(access1)
	if err != nil {
		t.Error(err)
		return
	}

	err = dto.Storage.Add(access2)
	if err != nil {
		t.Error(err)
		return
	}

	err = dto.Storage.Save(access3)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("find rows in DB")

	rows, err := dto.Storage.List()
	if err != nil {
		t.Error(err)
		return
	}

	if len(rows) != 3 {
		t.Errorf("rows length in db not equal 3 (%d)", len(rows))
		return
	}

	t.Log("remove one row")
	err = dto.Storage.Remove(access2)
	if err != nil {
		t.Error(err)
		return
	}

	rows, err = dto.Storage.List()
	if err != nil {
		t.Error(err)
		return
	}

	if len(rows) != 2 {
		t.Errorf("rows length in db not equal 2 (%d)", len(rows))
		return
	}

	t.Log("get one row by equal filter")
	rows, err = dto.Storage.List(params.NewEq("name", access3.Name()))
	if err != nil {
		t.Error(err)
		return
	}

	if len(rows) != 1 {
		t.Errorf("rows length in db not equal 1 (%d)", len(rows))
		return
	}

	if rows[0].Name() != access3.Name() {
		t.Errorf("incorrect row with name \"%s\" (ecpected \"%s\")", rows[0].Name(), access3.Name())
		return
	}

	t.Log("get one row by like filter")
	rows, err = dto.Storage.List(params.NewLike("name", access1.Name()+"%"))
	if err != nil {
		t.Error(err)
		return
	}

	if len(rows) != 1 {
		t.Errorf("rows length in db not equal 1 (%d)", len(rows))
		return
	}

	if rows[0].Name() != access1.Name() {
		t.Errorf("incorrect row with name \"%s\" (ecpected \"%s\")", rows[0].Name(), access1.Name())
		return
	}
}
