package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const (
	fileName = "config.ini"
)

type item struct {
	Name      string
	Issuer    string
	Key       string
	Algorithm string
	Digits    int
}

type Config struct {
	Items []item
}

func defaultConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "clotp")
}

// loadConfig parses `f` ini config to Config struct
func loadConfig(f io.Reader) (*Config, error) {
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

		s := item{
			Name:      section.Name(),
			Issuer:    section.Key("issuer").String(),
			Key:       key.String(),
			Algorithm: section.Key("algorithm").String(),
			Digits:    section.Key("digits").MustInt(),
		}

		config.Items = append(config.Items, s)
	}

	return &config, nil
}

// NewConfig reads config file if it exist or creates new empty one if doesn't
func NewConfig() (*Config, error) {
	dir := defaultConfigDir()
	cfgPath := filepath.Join(dir, fileName)

	if !fileExists(cfgPath) {
		if err := os.Mkdir(dir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		f, err := os.Create(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("failred to create config file: %w", err)
		}

		return loadConfig(f)
	}

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failied to read config file: %w", err)
	}

	return loadConfig(bytes.NewBuffer(b))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
