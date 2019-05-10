package api

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Org struct {
		Admin string `yaml:"admin"`
		Name  string `yaml:"name"`
	} `yaml:"org"`
	User struct {
		Name string `yaml:"name"`
	} `yaml:"user"`
	Ca struct{
		Host string
	}
}

func LoadConfiguration(file string) (*Config, error) {
	var config *Config
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to open configuration file")
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to parse configuration file JSON")
	}
	return config, nil
}