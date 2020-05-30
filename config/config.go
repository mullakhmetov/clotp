package config

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"os"
	"path/filepath"
)

const (
	defaultConfigName = "config.ini"
	defaultAlgorithm  = "sha1"
)

type parseAlgorithmFn func(string) (func() hash.Hash, error)

var supportedAlgorithms = []string{
	"sha1",
	"sha256",
	"sha512",
}

type Item struct {
	Name      string `ini:"-"`
	Issuer    string `ini:"issuer,omitempty"`
	Key       string `ini:"secret"`
	Algorithm string `ini:"algorithm,omitempty"`
	digest    func() hash.Hash
	Digits    int `ini:"digits,omitempty"`
	Step      int `ini:"step,omitempty"`
}

func (i Item) Validate() bool {
	if i.Name == "" {
		return false
	}

	if i.Key == "" {
		return false
	}

	if i.Step < 0 {
		return false
	}

	return true
}

func (i Item) Digest() func() hash.Hash {
	return i.digest
}

type Opts struct {
	path     string
	filename string
}

type mapper interface {
	Read() ([]*Item, error)
	Write(items []*Item) error
}

type Config struct {
	mapper           mapper
	parseAlgorithmFn parseAlgorithmFn
	itemNames        map[string]struct{}
	Items            []*Item
}

// Read reads config via mapper
func (c *Config) Read() error {
	items, err := c.mapper.Read()
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := c.add(item); err != nil {
			return err
		}
	}

	return nil
}

// Writes config via mapper
func (c Config) Write() error {
	return c.mapper.Write(c.Items)
}

// Add adds given item to config
func (c *Config) Add(item *Item) error {
	if ok := item.Validate(); !ok {
		return ErrInvalidItem
	}

	return c.add(item)
}

func (c *Config) add(item *Item) error {
	if c.itemNames == nil {
		c.itemNames = make(map[string]struct{})
	}

	if _, ok := c.itemNames[item.Name]; ok {
		return fmt.Errorf("%w: %s", ErrItemAlreadyExists, item.Name)
	}

	c.itemNames[item.Name] = struct{}{}
	c.Items = append(c.Items, item)

	return nil
}

// NewConfig reads config file if it exist or creates new empty one if doesn't
func NewConfig(opts Opts) (*Config, error) {
	if opts.path == "" {
		opts.path = defaultConfigDir()
	}

	if opts.filename == "" {
		opts.filename = defaultConfigName
	}

	cfg := &Config{
		mapper:           NewIniMapper(opts, parseAlgorithm),
		parseAlgorithmFn: parseAlgorithm,
	}

	if err := cfg.Read(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "clotp")
}

func parseAlgorithm(a string) (func() hash.Hash, error) {
	if a == "" {
		a = defaultAlgorithm
	}

	switch a {
	case "sha1":
		return sha1.New, nil
	case "sha256":
		return sha256.New, nil
	case "sha512":
		return sha512.New, nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", a)
	}
}
