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
	defaultConfigName = "config.ini"
)

type Item struct {
	Name      string `ini:"-"`
	Issuer    string `ini:"issuer,omitempty"`
	Key       string `ini:"secret"`
	Algorithm string `ini:"algorithm,omitempty"`
	Digits    int    `ini:"digits,omitempty"`
}

func (i Item) Validate() bool {
	if i.Name == "" {
		return false
	}

	if i.Key == "" {
		return false
	}

	return true
}

type Opts struct {
	path     string
	filename string
}

type Config struct {
	opts  Opts
	Items map[string]Item
}

// Read reads config from underlying file
func (c *Config) Read(r io.Reader) error {
	cfg, err := ini.Load(r)
	if err != nil {
		return err
	}

	for _, section := range cfg.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}

		item := Item{Name: section.Name()}
		if err := section.MapTo(&item); err != nil {
			return err
		}

		if item.Key == "" {
			continue
		}

		if err := c.add(item); err != nil {
			return err
		}
	}

	return nil
}

// Writes config to underlying file
func (c Config) Write(w io.WriteCloser) error {
	cfg := ini.Empty()

	for _, item := range c.Items {
		section, err := cfg.NewSection(item.Name)
		if err != nil {
			return err
		}
		err = section.ReflectFrom(&item)
		if err != nil {
			return err
		}
	}

	if _, err := cfg.WriteTo(w); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

// Add adds given item to config
func (c *Config) Add(item Item) error {
	if ok := item.Validate(); !ok {
		return ErrInvalidItem
	}

	return c.add(item)
}

func (c Config) List() []Item {
	items := make([]Item, 0, len(c.Items))
	for _, item := range c.Items {
		items = append(items, item)
	}

	return items
}

func (c *Config) add(item Item) error {
	if c.Items == nil {
		c.Items = make(map[string]Item)
	}

	if _, ok := c.Items[item.Name]; ok {
		return fmt.Errorf("%w: %s", ErrItemAlreadyExists, item.Name)
	}

	c.Items[item.Name] = item

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

	cfg := &Config{opts: opts}

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
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func defaultConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "clotp")
}
