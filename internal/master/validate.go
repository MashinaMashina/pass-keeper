package master

import (
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"pass-keeper/pkg/encrypt"
	"pass-keeper/pkg/filesystem"
)

func (m *master) validateConfig() error {
	masterFile, err := homedir.Expand(m.masterConfig.Get("file"))
	if err != nil {
		return err
	}

	exists, err := filesystem.Exists(masterFile)
	if err != nil {
		return err
	}

	if !exists {
		return m.masterFileDialog(masterFile, "Не найден файл с мастер паролем, введите мастер пароль:")
	}

	bytes, err := ioutil.ReadFile(masterFile)
	if err != nil {
		return err
	}

	deviceKey, err := m.DeviceCryptoKey()
	if err != nil {
		return err
	}

	pwd, err := encrypt.DecryptAES(deviceKey, string(bytes))
	if err != nil {
		return m.masterFileDialog(masterFile, "Не верный мастер пароль, введите мастер пароль заново:")
	}

	_, err = hex.DecodeString(pwd)
	if err != nil {
		return m.masterFileDialog(masterFile, "Не верный мастер пароль, введите мастер пароль заново:")
	}

	m.masterConfig.Set("password", pwd)

	return nil
}

func (m *master) masterFileDialog(masterFile, message string) error {
	fmt.Println(message)
	var pwd string
	_, err := fmt.Scanln(&pwd)
	if err != nil {
		return err
	}

	return m.saveMasterPassword(masterFile, pwd)
}

func (m *master) saveMasterPassword(masterFile, pwd string) error {
	pwd = m.Hash(pwd)

	m.masterConfig.Set("password", pwd)

	deviceKey, err := m.DeviceCryptoKey()
	if err != nil {
		return err
	}

	encoded, err := encrypt.EncryptAES(deviceKey, pwd)
	if err != nil {
		return errors.Wrap(err, "encoding master password")
	}

	err = ioutil.WriteFile(masterFile, []byte(encoded), 0777)
	if err != nil {
		return errors.Wrap(err, "write master password failed")
	}

	return nil
}