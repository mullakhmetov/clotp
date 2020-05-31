package main

import (
	"crypto/sha1" //nolint:gosec // used in hmac only, see RFC 4226 B.2. section
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/ini.v1"
)

var (
	fullConfig = []byte(`[Name-1]
issuer=issuer-1
secret=secret-key-1
algorithm=sha1
digits=6
step=30
`)

	onlySecret = []byte(`[Name-2]
secret=secret-key-2
`)

	noSecret = []byte(`[Name-3]
issuer=issuer-3
algorithm=sha1
digits=6
`)

	multiple = []byte(`[Name-4]
issuer=issuer-4
secret=secret-key-4
algorithm=sha1
digits=6
step=30

[Name-5]
issuer=issuer-5
secret=secret-key-5
algorithm=sha1
digits=6
step=60
`)
)

func parse(string) (func() hash.Hash, error) {
	return sha1.New, nil
}

func TestInitMapper_Read_Create(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	mapper := NewIniMapper(Opts{path: dir}, parse)

	items, err := mapper.Read()
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(items, []*Item{}) {
		t.Errorf("items should be nil, got: %+v", items)
	}

	cfgPath := filepath.Join(dir, mapper.opts.filename)
	_, err = os.Stat(cfgPath)
	if err != nil {
		panic(err)
	}
}

func TestMapper_Read_Open(t *testing.T) {
	ini.PrettyEqual = false
	ini.PrettyFormat = false

	for _, c := range []struct {
		name   string
		mapper mapper
		input  []byte
		want   []*Item
	}{
		{
			name:  "check fields parsing",
			input: fullConfig,
			want: []*Item{
				{Name: "Name-1", Issuer: "issuer-1", Key: "secret-key-1", Algorithm: "sha1", Digits: 6, Step: 30},
			},
		},
		{
			name:  "check empty fields",
			input: onlySecret,
			want: []*Item{
				{Name: "Name-2", Key: "secret-key-2"},
			},
		},

		{
			name:  "check required secret field",
			input: noSecret,
			want:  []*Item{},
		},
		{
			name:  "check multiple secrets",
			input: multiple,
			want: []*Item{
				{Name: "Name-4", Issuer: "issuer-4", Key: "secret-key-4", Algorithm: "sha1", Digits: 6, Step: 30},
				{Name: "Name-5", Issuer: "issuer-5", Key: "secret-key-5", Algorithm: "sha1", Digits: 6, Step: 60},
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "prefix-"+c.name)
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(dir)

			file, err := ioutil.TempFile(dir, "*"+defaultConfigName)
			if err != nil {
				panic(err)
			}
			defer os.Remove(file.Name())

			if _, err := file.Write(c.input); err != nil {
				panic(err)
			}

			mapper := NewIniMapper(Opts{path: dir, filename: filepath.Base(file.Name())}, parse)

			items, err := mapper.Read()
			if err != nil {
				panic(err)
			}

			// function type is incomparable
			for _, i := range items {
				i.digest = nil
			}

			if !reflect.DeepEqual(c.want, items) {
				t.Errorf("wrong config, want: %+v != got: %+v", c.want, items)
			}
		})
	}
}
