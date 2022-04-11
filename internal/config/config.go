package config

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type Config struct {
	values    map[string]interface{}
	virtuals  []string
	initFile  string
	hasChange bool
}

func NewConfig() *Config {
	return &Config{
		values: map[string]interface{}{},
	}
}

func (c *Config) HasChange() bool {
	return c.hasChange
}

func (c *Config) Set(k string, v interface{}) {
	if c.Get(k) == v {
		return
	}

	isVirt := false
	for _, val := range c.virtuals {
		if val == k {
			isVirt = true
			break
		}
	}

	if !isVirt {
		c.hasChange = true
	}

	c.values[k] = v
}

func (c *Config) Virtual() []string {
	return c.virtuals
}

func (c *Config) SetVirtual(v []string) {
	c.virtuals = v
}

func (c *Config) Get(k string) interface{} {
	if val, exists := c.values[k]; exists {
		return val
	}

	return nil
}

func (c *Config) String(k string) string {
	val := c.Get(k)
	switch val.(type) {
	case float64:
		return fmt.Sprintf("%v", val)
	case string:
		return val.(string)
	}

	return ""
}

func (c *Config) Map(k string) map[string]string {
	val := c.Get(k)
	switch val.(type) {
	case map[string]string:
		return val.(map[string]string)
	case map[string]interface{}:
		res := make(map[string]string)
		for k, v := range val.(map[string]interface{}) {
			res[k] = fmt.Sprint(v)
		}

		return res
	}

	return nil
}

func (c *Config) Slice(k string) []string {
	val := c.Get(k)
	switch val.(type) {
	case []string:
		return val.([]string)
	case []interface{}:
		var res []string
		for _, v := range val.([]interface{}) {
			res = append(res, fmt.Sprint(v))
		}

		return res
	}

	return nil
}

func (c *Config) InitFromFile(file string) error {
	var err error
	file, err = homedir.Expand(file)

	if err != nil {
		return errors.Wrap(err, "expand "+file)
	}

	c.initFile = file

	bytes, err := ioutil.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	return c.InitFromData(bytes)
}

func (c *Config) InitFromData(b []byte) error {
	err := json.Unmarshal(b, &c.values)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SaveToFile(file ...string) error {
	if !c.hasChange {
		return nil
	}

	var f string
	if len(file) == 0 {
		f = c.initFile
	} else {
		f = file[0]
	}

	bytes, err := c.SaveToData()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SaveToData() ([]byte, error) {
	values := c.values
	for _, k := range c.virtuals {
		delete(values, k)
	}

	bytes, err := json.MarshalIndent(values, "", "	")
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
