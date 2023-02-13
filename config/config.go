package config

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Telegram struct {
		Token   string `yaml:"token"`
		Channel int64  `yaml:"channel"`
	} `yaml:"telegram"`
	Email struct {
		Host          string `yaml:"host"`
		Port          uint32 `yaml:"port"`
		User          string `yaml:"user"`
		TemplatesPath string `yaml:"templates_path"`
		Password      string `yaml:"-"`
	} `yaml:"email"`
	RedPanda struct {
		Host string `yaml:"host"`
	} `yaml:"red_panda"`
	Service struct {
		Port    string `yaml:"port"`
		Address string `yaml:"address"`
	}
}

func Get(path string) (*Config, error) {
	cfg := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(file)

	decoder := yaml.NewDecoder(buf)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	cfg.Email.Password = os.Getenv("EMAIL_PWD")

	return cfg, nil
}
