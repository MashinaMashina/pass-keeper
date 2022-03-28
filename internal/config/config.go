package config

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"pass-keeper/pkg/filesystem"
)

var FileNotExists = errors.New("config file not exists")

type Config struct {
	values map[string]*Part
	file   string
}

func New() (*Config, error) {
	cfg := Config{
		values: map[string]*Part{},
	}

	var err error
	cfg.file, err = homedir.Expand("~/.pass-keeper.json")

	if err != nil {
		return nil, errors.Wrap(err, "expand ~/.pass-keeper.json")
	}

	return &cfg, nil
}

func (c *Config) GenFile() error {
	for key, handler := range c.values {
		for _, field := range handler.InstallFields() {
			name := fmt.Sprintf("%s.%s", key, field)

			if humanName, exists := handler.FieldNames()[field]; exists {
				name = humanName
			}

			text := fmt.Sprintf("Введите %s", name)

			def, exists := handler.Default(field)
			if exists {
				text = fmt.Sprintf("%s (по-умолчанию \"%s\")", text, def)
			}

			text = fmt.Sprintf("%s:", text)

			fmt.Println(text)

			var val string
			_, err := fmt.Scanln(&val)
			if err != nil && err.Error() != "unexpected newline" {
				return err
			}

			if val == "" {
				val = def
			}

			c.values[key].Set(field, val)
		}
	}

	return c.Save()
}

func (c *Config) Save() error {
	bytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.file, bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) AddPart(name string, handler *Part) error {
	c.values[name] = handler

	return nil
}

func (c *Config) Part(name string) *Part {
	return c.values[name]
}

func (c *Config) Init() error {
	for _, handler := range c.values {
		if err := handler.Init(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) LoadUserConfig() error {
	err := c.Load(c.file)

	if errors.Is(err, FileNotExists) {
		err = c.GenFile()
		if err != nil {
			return errors.Wrap(err, "creating config")
		}
	} else if err != nil {
		return errors.Wrap(err, "read config")
	}

	return nil
}

func (c *Config) Load(filename string) error {
	exists, err := filesystem.Exists(filename)
	if err != nil {
		return err
	}

	if !exists {
		return FileNotExists
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.values)
}

func (c *Config) UnmarshalJSON(bytes []byte) error {
	var m map[string]map[string]string

	err := json.Unmarshal(bytes, &m)

	if err != nil {
		return err
	}

	for key, value := range m {
		if _, exists := c.values[key]; exists {
			c.values[key].Load(value)
		}
	}

	return nil
}
