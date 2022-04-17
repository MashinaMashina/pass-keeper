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

func (m *master) initConfig() error {
	if m.Config.String("master.password") != "" {
		return nil
	}

	if m.Config.String("master.file") == "" {
		m.Config.Set("master.file", "~/.pass-keeper.master")
	}

	masterFile, err := homedir.Expand(m.Config.String("master.file"))
	if err != nil {
		return err
	}

	exists, err := filesystem.Exists(masterFile)
	if err != nil {
		return err
	}

	if !exists {
		return m.masterFileDialog(masterFile, "Master password file not found, please enter master password:")
	}

	bytes, err := ioutil.ReadFile(masterFile)
	if err != nil {
		return err
	}

	deviceKey, err := m.deviceCryptoKey(m.Config.String("main.key"))
	if err != nil {
		return err
	}

	pwd, err := encrypt.DecryptAES(deviceKey, string(bytes))
	if err != nil {
		return m.masterFileDialog(masterFile, "Wrong master password, please re-enter the master password:")
	}

	_, err = hex.DecodeString(pwd)
	if err != nil {
		return m.masterFileDialog(masterFile, "Wrong master password, please re-enter the master password:")
	}

	m.Config.Set("master.password", pwd)

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
	pwd = m.hash(pwd)

	m.Config.Set("master.password", pwd)

	deviceKey, err := m.deviceCryptoKey(m.Config.String("main.key"))
	if err != nil {
		return err
	}

	encoded, err := encrypt.EncryptAES(deviceKey, pwd)
	if err != nil {
		return errors.Wrap(err, "encoding master password")
	}

	err = ioutil.WriteFile(masterFile, []byte(encoded), 0644)
	if err != nil {
		return errors.Wrap(err, "write master password failed")
	}

	return nil
}
