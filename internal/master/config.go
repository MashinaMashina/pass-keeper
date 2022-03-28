package master

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"pass-keeper/pkg/filesystem"
)

var secret = "e21r5naewwps5yx2023yl24fnl1gos5s"

func (m *master) fillConfig() {
	t := m.masterConfig.TemporaryFields()
	t = append(t, "password")
	m.masterConfig.SetTemporaryFields(t)

	d := m.masterConfig.DefaultValues()
	d["file"] = "~/.pass-keeper.master"
	m.masterConfig.SetDefaultValues(d)

	i := m.masterConfig.InstallFields()
	i = append(i, "file")
	m.masterConfig.SetInstallFields(i)

	f := m.masterConfig.FieldNames()
	f["file"] = "файл с мастер паролем"
	m.masterConfig.SetFieldNames(f)

	m.masterConfig.SetInit(m.validateConfig)
}

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
		return m.masterFileDialog(masterFile)
	}

	bytes, err := ioutil.ReadFile(masterFile)
	if err != nil {
		return err
	}

	m.masterConfig.Set("password", string(bytes))

	return nil
}

func (m *master) masterFileDialog(masterFile string) error {
	fmt.Println("Не найден файл с мастер паролем, введите мастер пароль:")
	var pwd string
	_, err := fmt.Scanln(&pwd)
	if err != nil {
		return err
	}

	pwd = m.Encode(pwd)

	m.masterConfig.Set("password", pwd)

	err = ioutil.WriteFile(masterFile, []byte(pwd), 0777)
	if err != nil {
		return errors.Wrap(err, "write master password failed")
	}

	fmt.Println("Lalala:", m.masterConfig.Get("password"))

	return nil
}

func (m *master) Encode(s string) string {
	bytes := md5.Sum([]byte(s + secret))

	return hex.EncodeToString(bytes[:])
}
