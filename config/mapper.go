package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

func NewIniMapper(opts Opts, fn parseAlgorithmFn) *IniMapper {
	path := filepath.Join(opts.path, opts.filename)
	return &IniMapper{opts, path, fn}
}

type IniMapper struct {
	opts Opts
	path string

	parseAlgorithmFn
}

// Read reads config items from reader
func (m *IniMapper) Read() ([]*Item, error) {
	r, err := m.reader()
	if err != nil {
		return nil, err
	}

	cfg, err := ini.Load(r)
	if err != nil {
		return nil, err
	}

	items := make([]*Item, 0)

	for _, section := range cfg.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}

		item := &Item{Name: section.Name()}
		if err := section.MapTo(&item); err != nil {
			return nil, err
		}

		d, err := m.parseAlgorithmFn(item.Algorithm)
		if err != nil {
			return nil, err
		}

		item.digest = d

		if item.Key == "" {
			continue
		}

		items = append(items, item)
	}

	return items, nil
}

// Write writes config items to writer
func (m IniMapper) Write(items []*Item) error {
	w, err := m.writer()
	if err != nil {
		return err
	}

	cfg := ini.Empty()

	for _, item := range items {
		section, err := cfg.NewSection(item.Name)
		if err != nil {
			w.Close()
			return err
		}

		err = section.ReflectFrom(item)
		if err != nil {
			w.Close()
			return err
		}
	}

	if _, err := cfg.WriteTo(w); err != nil {
		w.Close()
		return err
	}

	return w.Close()
}

func (m IniMapper) reader() (io.Reader, error) {
	return m.readWriterCloser()
}

func (m IniMapper) writer() (io.WriteCloser, error) {
	return m.readWriterCloser()
}

func (m IniMapper) readWriterCloser() (io.ReadWriteCloser, error) {
	if m.opts.path == "" {
		m.opts.path = defaultConfigDir()
	}

	if m.opts.filename == "" {
		m.opts.filename = defaultConfigName
	}

	if !pathExists(m.opts.path) {
		if err := os.Mkdir(m.opts.path, 0700); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	cfgPath := filepath.Join(m.opts.path, m.opts.filename)

	var getFile func(string) (*os.File, error)

	if !pathExists(cfgPath) {
		getFile = os.Create
	} else {
		getFile = func(p string) (*os.File, error) { return os.OpenFile(p, os.O_RDWR, 0) }
	}

	f, err := getFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failied to open config file: %w", err)
	}

	return f, err
}

func pathExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
