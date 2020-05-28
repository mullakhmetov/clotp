package config

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"

	"gopkg.in/ini.v1"
)

func NewIniMapper() *IniMapper {
	return &IniMapper{}
}

type IniMapper struct{}

// Read reads config items from reader
func (m *IniMapper) Read(r io.Reader) ([]*Item, error) {
	var items []*Item

	cfg, err := ini.Load(r)
	if err != nil {
		return nil, err
	}

	for _, section := range cfg.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}

		item := &Item{Name: section.Name()}
		if err := section.MapTo(&item); err != nil {
			return nil, err
		}

		d, err := parseAlgorithm(item.Algorithm)
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
func (m IniMapper) Write(items []*Item, w io.WriteCloser) error {
	cfg := ini.Empty()

	for _, item := range items {
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

func parseAlgorithm(a string) (func() hash.Hash, error) {
	if a == "" {
		a = "sha1"
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
