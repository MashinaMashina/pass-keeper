package master

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/denisbrodbeck/machineid"
)

func (m *master) DeviceCryptoKey() ([]byte, error) {
	machineId, err := machineid.ID()

	if err != nil {
		return nil, err
	}

	bytes := md5.Sum([]byte(machineId + m.config.Part("main").Get("key")))

	return bytes[:], nil
}

func (m *master) Hash(s string) string {
	bytes := md5.Sum([]byte(s))

	return hex.EncodeToString(bytes[:])
}
