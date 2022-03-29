package encrypt

import (
	"encoding/hex"
	"testing"
)

func TestKey16(t *testing.T) {
	key, err := hex.DecodeString("c4ca4238a0b923820dcc509a6f75849b")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Key length:", len(key))

	str := "some string some string"

	encoded, err := EncryptAES(key, str)
	if err != nil {
		t.Error(err)
		return
	}

	decoded, err := DecryptAES(key, encoded)
	if err != nil {
		t.Error(err)
		return
	}

	if decoded != str {
		t.Errorf("DecryptAES(EncryptAES(%s)) != %s", str, decoded)
	}
}

func TestKey32(t *testing.T) {
	key, err := hex.DecodeString("c4ca4238a0b923820dcc509a6f75849b" +
		"c4ca4238a0b923820dcc509a6f75849b")

	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Key length:", len(key))

	str := "some string some string"

	encoded, err := EncryptAES(key, str)
	if err != nil {
		t.Error(err)
		return
	}

	decoded, err := DecryptAES(key, encoded)
	if err != nil {
		t.Error(err)
		return
	}

	if decoded != str {
		t.Errorf("DecryptAES(EncryptAES(%s)) != %s", str, decoded)
	}
}
