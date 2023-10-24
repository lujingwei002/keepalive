package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	c *Config
)

type Server struct {
	Heartbeat time.Duration `yaml:"heartbeat"`
	Type      string        `yaml:"type"`
	Ws        *struct {
		Bind string `yaml:"bind"`
		Path string `yaml:"path"`
	} `yaml:"ws"`
}

type Backend struct {
	// push to backend config
	Sink struct {
		Type string `yaml:"type"`

		Redis *struct {
			Host string  `yaml:"string"`
			Port int     `yaml:"port"`
			Key  string  `yaml:"key"`
			Db   int     `yaml:"0"`
			Auth *string `yaml:"auth"`
		} `yaml:"redis"`

		Grpc *struct {
			Address string `yaml:"address"`
		} `yaml:"grpc"`
	} `yaml:"sink"`

	// pull from backend config
	Source struct {
		Type string `yaml:"type"`
	} `yaml:"source"`
}

type Config struct {
	Servers  map[string]Server  `yaml:"server"`
	Backends map[string]Backend `yaml:"backend"`
}

func init() {
	c = &Config{}
}

func Parse(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(content, &c); err != nil {
		return err
	}
	fmt.Printf("eee %+v\n", c)
	return nil
}

func Get() *Config {
	return c
}
