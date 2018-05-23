package model

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Config *config

type config struct {
	Listen    string `yaml:"listen"`
	LogLevel  string `yaml:"log_level"`
	Initiator string `yaml:"initiator"`
	Secret    string `yaml:"secret"`
	Database  string `yaml:"databse"`
}

//NewConfig returns config with default values.
func NewConfig() *config {
	return &config{
		Listen:    ":8080",
		LogLevel:  "debug",
		Initiator: "admin@systemli.org",
		Secret:    "secret",
		Database:  "ticker.db",
	}
}

//LoadConfig loads config from file.
func LoadConfig(path string) *config {
	c := NewConfig()

	yml, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithField("path", path).Panic("failed to open config")
	}
	err = yaml.Unmarshal(yml, &c)
	if err != nil {
		log.Panic("failed to parse config")
	}

	return c
}
