package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/ini.v1"
)

var (
	fullConfig = []byte(`[Name-1]
issuer=issuer-1
secret=secret-key-1
algorithm=sha1
digits=6
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

[Name-5]
issuer=issuer-5
secret=secret-key-5
algorithm=sha1
digits=6
`)
)

func TestConfig_Read(t *testing.T) {
	ini.PrettyEqual = false
	ini.PrettyFormat = false

	for _, c := range []struct {
		name  string
		input []byte
		want  *Config
	}{
		{
			name:  "check fields parsing",
			input: fullConfig,
			want: &Config{Items: []item{
				{Name: "Name-1", Issuer: "issuer-1", Key: "secret-key-1", Algorithm: "sha1", Digits: 6},
			}},
		},
		{
			name:  "check empty fields",
			input: onlySecret,
			want: &Config{Items: []item{
				{Name: "Name-2", Key: "secret-key-2"}},
			},
		},
		{
			name:  "check required secret field",
			input: noSecret,
			want:  &Config{},
		},
		{
			name:  "check multiple secrets",
			input: multiple,
			want: &Config{Items: []item{
				{Name: "Name-4", Issuer: "issuer-4", Key: "secret-key-4", Algorithm: "sha1", Digits: 6},
				{Name: "Name-5", Issuer: "issuer-5", Key: "secret-key-5", Algorithm: "sha1", Digits: 6},
			}},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			cfg := &Config{}
			if err := cfg.Read(bytes.NewReader(c.input)); err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(c.want, cfg) {
				t.Errorf("wrong config loading, want: %+v != got: %+v", c.want, cfg)
			}
		})
	}
}

type stubWriteCloser struct {
	*bytes.Buffer
	closed bool
}

func (wc *stubWriteCloser) Close() error {
	wc.closed = true
	return nil
}

func TestConfig_Write(t *testing.T) {
	ini.PrettyEqual = false
	ini.PrettyFormat = false

	for _, c := range []struct {
		name string
		cfg  *Config
		want []byte
	}{
		{
			name: "check fields parsing",
			cfg: &Config{Items: []item{
				{Name: "Name-1", Issuer: "issuer-1", Key: "secret-key-1", Algorithm: "sha1", Digits: 6},
			}},
			want: fullConfig,
		},
		{
			name: "check empty fields",
			cfg: &Config{Items: []item{
				{Name: "Name-2", Key: "secret-key-2"}},
			},
			want: onlySecret,
		},
		{
			name: "check multiple secrets",
			cfg: &Config{Items: []item{
				{Name: "Name-4", Issuer: "issuer-4", Key: "secret-key-4", Algorithm: "sha1", Digits: 6},
				{Name: "Name-5", Issuer: "issuer-5", Key: "secret-key-5", Algorithm: "sha1", Digits: 6},
			}},
			want: multiple,
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			wc := &stubWriteCloser{bytes.NewBufferString(""), false}
			if err := c.cfg.Write(wc); err != nil {
				panic(err)
			}

			if !wc.closed {
				t.Error("config not closed")
			}

			if got := strings.Trim(wc.String(), "\n"); got != strings.Trim(string(c.want), "\n") {
				t.Errorf("wrong config writing, want: %s != got: %s", c.want, got)
			}
		})
	}
}

func TestNewConfig_Init(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	got, err := NewConfig(Opts{path: dir})
	if err != nil {
		panic(err)
	}

	want := &Config{opts: Opts{path: dir, filename: defaultConfigName}}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("wrong config, want: %+v != got: %+v", want, got)
	}
}

func TestNewConfig_Read(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	file, err := ioutil.TempFile(dir, "*"+defaultConfigName)
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write(multiple); err != nil {
		panic(err)
	}

	got, err := NewConfig(Opts{path: dir, filename: filepath.Base(file.Name())})
	if err != nil {
		fmt.Println(file.Name())
		panic(err)
	}

	want := &Config{
		Items: []item{
			{Name: "Name-4", Issuer: "issuer-4", Key: "secret-key-4", Algorithm: "sha1", Digits: 6},
			{Name: "Name-5", Issuer: "issuer-5", Key: "secret-key-5", Algorithm: "sha1", Digits: 6},
		},
		opts: Opts{path: dir, filename: filepath.Base(file.Name())},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wrong config, want: %+v != got: %+v", want, got)
	}
}
