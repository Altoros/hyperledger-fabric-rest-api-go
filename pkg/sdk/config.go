package sdk

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Org struct {
		Admin string
		Name  string
	}
	User struct {
		Name string
	}
	Ca struct {
		Url         string
		Tls         bool
		TlsCertFile string `yaml:"tlsCertFile"`
		Address     string
		Protocol    string
	}
	Orderer ApiOrderer
	Peers map[string]ApiPeer
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

	config.Ca.Protocol = "http"
	if config.Ca.Tls {
		config.Ca.Protocol = "https"
	}

	config.Ca.Address = fmt.Sprintf("%s://%s", config.Ca.Protocol, config.Ca.Url)

	return config, nil
}
