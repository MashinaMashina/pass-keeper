package test

import (
	"encoding/hex"
	"pass-keeper/internal/master"
	"pass-keeper/pkg/encrypt"
	"testing"
)

func TestHash(t *testing.T) {
	dto, err := testingDTO()
	if err != nil {
		t.Error(err)
	}
	defer dto.Storage.Close()

	m := master.New(dto)

	hashed := m.Hash("test")

	_, err = hex.DecodeString(hashed)
	if err != nil {
		t.Errorf("hash must be hex decodable: %v", err)
	}
}

func TestDeviceCryptoKey(t *testing.T) {
	dto, err := testingDTO()
	if err != nil {
		t.Error(err)
	}
	defer dto.Storage.Close()

	m := master.New(dto)

	key, err := m.DeviceCryptoKey()
	if err != nil {
		t.Error(err)
	}

	key = append(key, key...)

	t.Log("Key len", len(key))

	str := "some string somesome string some"
	result, err := encrypt.EncryptAES(key, str)
	if err != nil {
		t.Error(err)
	}

	decrypt, err := encrypt.DecryptAES(key, result)
	if err != nil {
		t.Error(err)
	}

	if str != decrypt {
		t.Errorf("DecryptAES(EncryptAES(%s)) != %s", str, decrypt)
	}
}
