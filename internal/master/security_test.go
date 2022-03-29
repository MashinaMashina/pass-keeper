package master

import (
	"encoding/hex"
	"pass-keeper/pkg/encrypt"
	"testing"
)

func TestHash(t *testing.T) {
	cfg, storage, err := configAndStorage()
	if err != nil {
		t.Error(err)
	}
	defer storage.Close()

	m := New(storage, cfg)

	hashed := m.Hash("test")

	_, err = hex.DecodeString(hashed)
	if err != nil {
		t.Errorf("hash must be hex decodable: %v", err)
	}
}

func TestDeviceCryptoKey(t *testing.T) {
	cfg, storage, err := configAndStorage()
	if err != nil {
		t.Error(err)
	}
	defer storage.Close()

	m := New(storage, cfg)

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
