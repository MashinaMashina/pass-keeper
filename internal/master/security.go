package master

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/denisbrodbeck/machineid"
)

func (m *master) deviceCryptoKey(mainKey string) ([]byte, error) {
	machineId, err := machineid.ID()

	if err != nil {
		return nil, err
	}

	bytes := md5.Sum([]byte(machineId + mainKey))

	return bytes[:], nil
}

func (m *master) hash(s string) string {
	bytes := md5.Sum([]byte(s))

	return hex.EncodeToString(bytes[:])
}
