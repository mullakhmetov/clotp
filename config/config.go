package config

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const (
	fileName = "config.ini"
)

type secret struct {
	Name      string
	Issuer    string
	Key       string
	Algorithm string
	Digits    int
}

type Config struct {
	Secrets []secret
}

func DefaultConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "clotp")
}

func LoadConfig(f io.Reader) (*Config, error) {
	var config Config

	cfg, err := ini.Load(f)
	if err != nil {
		return nil, err
	}

	for _, section := range cfg.Sections() {
		key, err := section.GetKey("secret")
		if err != nil {
			continue
		}

		s := secret{
			Name:      section.Name(),
			Issuer:    section.Key("issuer").String(),
			Key:       key.String(),
			Algorithm: section.Key("algorithm").String(),
			Digits:    section.Key("digits").MustInt(),
		}

		config.Secrets = append(config.Secrets, s)
	}

	return &config, nil
}
