package config

import (
	"io/ioutil"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"gopkg.in/yaml.v2"
)

// Conf config
var Conf *Config

// DefaultPath is default config file path
var DefaultPath = "./configs/kms.yaml"

// Config presents config
type Config struct {
	Port        string        `yaml:"port"`
	Model       string        `yaml:"model"`
	Logger      logger.Config `yaml:"log"`
	InternalNet client.Config `yaml:"internalNet"`
	Mysql       mysql2.Config `yaml:"mysql"`
	Redis       redis2.Config `yaml:"redis"`
	OAuth2Host  string        `yaml:"oauth2Host"`
}

// NewConfig load config file
func NewConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultPath
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return nil, err
	}

	return Conf, nil
}
