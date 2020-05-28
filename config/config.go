package config

import (
	"bytes"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultConfigName = "config.ini"
)

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

type Opts struct {
	path     string
	filename string
}

type mapper interface {
	Read(io.Reader) ([]*Item, error)
	Write(items []*Item, w io.WriteCloser) error
}

type Config struct {
	mapper    mapper
	opts      Opts
	itemNames map[string]struct{}
	Items     []*Item
}

// Read reads config via mapper
func (c *Config) Read(r io.Reader) error {
	items, err := c.mapper.Read(r)
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
func (c Config) Write(w io.WriteCloser) error {
	return c.mapper.Write(c.Items, w)
}

// Add adds given item to config
func (c *Config) Add(item *Item) error {
	if ok := item.Validate(); !ok {
		return ErrInvalidItem
	}

	return c.add(item)
}

func (c Config) List() []*Item {
	return c.Items
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

// NewDefaultConfig reads config file if it exist or creates new empty one if doesn't
func NewDefaultConfig() (*Config, error) {
	opts := Opts{
		path:     defaultConfigDir(),
		filename: defaultConfigName,
	}

	return NewConfig(opts)
}

// NewConfig reads config file if it exist or creates new empty one if doesn't
func NewConfig(opts Opts) (*Config, error) {
	if opts.path == "" {
		opts.path = defaultConfigDir()
	}

	if opts.filename == "" {
		opts.filename = defaultConfigName
	}

	if !pathExists(opts.path) {
		if err := os.Mkdir(opts.path, 0700); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	cfg := &Config{mapper: NewIniMapper(), opts: opts}

	cfgPath := filepath.Join(opts.path, opts.filename)

	if !pathExists(cfgPath) {
		f, err := os.Create(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("failred to create config file: %w", err)
		}

		if err := cfg.Read(f); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failied to read config file: %w", err)
	}

	if err := cfg.Read(bytes.NewBuffer(b)); err != nil {
		return nil, err
	}

	return cfg, nil
}

func pathExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func defaultConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "clotp")
}
