package config

import (
	"bytes"
	"errors"
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

func TestConfig_Read(t *testing.T) {
	ini.PrettyEqual = false
	ini.PrettyFormat = false

	for _, c := range []struct {
		name   string
		mapper mapper
		input  []byte
		want   []*Item
	}{
		{
			name:   "check fields parsing",
			mapper: NewIniMapper(),
			input:  fullConfig,
			want: []*Item{
				{Name: "Name-1", Issuer: "issuer-1", Key: "secret-key-1", Algorithm: "sha1", Digits: 6, Step: 30},
			},
		},
		{
			name:   "check empty fields",
			mapper: NewIniMapper(),
			input:  onlySecret,
			want: []*Item{
				{Name: "Name-2", Key: "secret-key-2"},
			},
		},

		{
			name:   "check required secret field",
			mapper: NewIniMapper(),
			input:  noSecret,
		},
		{
			name:   "check multiple secrets",
			mapper: NewIniMapper(),
			input:  multiple,
			want: []*Item{
				{Name: "Name-4", Issuer: "issuer-4", Key: "secret-key-4", Algorithm: "sha1", Digits: 6, Step: 30},
				{Name: "Name-5", Issuer: "issuer-5", Key: "secret-key-5", Algorithm: "sha1", Digits: 6, Step: 60},
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			cfg := &Config{mapper: c.mapper}
			if err := cfg.Read(bytes.NewReader(c.input)); err != nil {
				panic(err)
			}

			// function type is incomparable
			for _, i := range cfg.Items {
				i.digest = nil
			}

			if !reflect.DeepEqual(c.want, cfg.Items) {
				t.Errorf("wrong config loading, want: %+v != got: %+v", c.want, cfg.Items)
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
			cfg: &Config{
				mapper: NewIniMapper(),
				Items: []*Item{
					{Name: "Name-1", Issuer: "issuer-1", Key: "secret-key-1", Algorithm: "sha1", Digits: 6, Step: 30},
				},
			},
			want: fullConfig,
		},
		{
			name: "check empty fields",
			cfg: &Config{
				mapper: NewIniMapper(),
				Items: []*Item{
					{Name: "Name-2", Key: "secret-key-2"},
				},
			},
			want: onlySecret,
		},
		{
			name: "check multiple secrets",
			cfg: &Config{
				mapper: NewIniMapper(),
				Items: []*Item{
					{Name: "Name-4", Issuer: "issuer-4", Key: "secret-key-4", Algorithm: "sha1", Digits: 6, Step: 30},
					{Name: "Name-5", Issuer: "issuer-5", Key: "secret-key-5", Algorithm: "sha1", Digits: 6, Step: 60},
				},
			},
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

	want := &Config{mapper: NewIniMapper(), opts: Opts{path: dir, filename: defaultConfigName}}
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

	// function type is incomparable
	for _, i := range got.Items {
		i.digest = nil
	}

	want := &Config{
		mapper:    NewIniMapper(),
		itemNames: map[string]struct{}{"Name-4": {}, "Name-5": {}},
		Items: []*Item{
			{Name: "Name-4", Issuer: "issuer-4", Key: "secret-key-4", Algorithm: "sha1", Digits: 6, Step: 30},
			{Name: "Name-5", Issuer: "issuer-5", Key: "secret-key-5", Algorithm: "sha1", Digits: 6, Step: 60},
		},
		opts: Opts{path: dir, filename: filepath.Base(file.Name())},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wrong config, want: %+v != got: %+v", want, got)
	}
}

func TestItemValidate(t *testing.T) {
	for _, c := range []struct {
		name string
		item Item
		want bool
	}{
		{
			name: "no name",
			item: Item{Key: "1"},
			want: false,
		},
		{
			name: "no key",
			item: Item{Name: "n"},
			want: false,
		},
		{
			name: "negative step",
			item: Item{Name: "n", Step: -1},
			want: false,
		},
		{
			name: "valid",
			item: Item{Name: "n", Key: "k"},
			want: true,
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := c.item.Validate(); got != c.want {
				t.Errorf("wrong validation result, should be %t", c.want)
			}
		})
	}
}

func TestConfigAdd(t *testing.T) {
	config := Config{}

	for _, c := range []struct {
		name string
		item *Item
		err  error
	}{
		{
			name: "no name",
			item: &Item{Key: "1"},
			err:  ErrInvalidItem,
		},
		{
			name: "no key",
			item: &Item{Name: "n"},
			err:  ErrInvalidItem,
		},
		{
			name: "negative step",
			item: &Item{Name: "foo", Key: "bar", Step: -1},
			err:  ErrInvalidItem,
		},
		{
			name: "valid #1",
			item: &Item{Name: "n", Key: "k"},
		},
		{
			name: "duplicate",
			item: &Item{Name: "n", Key: "k2"},
			err:  ErrItemAlreadyExists,
		},
		{
			name: "valid #2",
			item: &Item{Name: "n2", Key: "k"},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			err := config.Add(c.item)
			if err != nil {
				if c.err == nil {
					t.Errorf("unwanted error: %v", err)
					return
				}

				if !errors.Is(err, c.err) {
					t.Errorf("error should match with %v", c.err)
					return
				}

				return
			}

			if c.err != nil {
				t.Error("error shouldn't be nil")
				return
			}
		})
	}

	want := []*Item{
		{Name: "n", Key: "k"},
		{Name: "n2", Key: "k"},
	}

	// function type is incomparable
	for _, i := range config.Items {
		i.digest = nil
	}

	if !reflect.DeepEqual(want, config.Items) {
		t.Errorf("wrong config items state, want: %+v != got: %+v", want, config.Items)
	}
}

func TestConfigList(t *testing.T) {
	config := Config{}

	items := []*Item{{Name: "n", Key: "k"}, {Name: "n2", Key: "k"}}
	for _, item := range items {
		if err := config.Add(item); err != nil {
			t.Errorf("unwanted error: %w", err)
			return
		}
	}

	if got := config.List(); !reflect.DeepEqual(items, got) {
		t.Errorf("wong config.List() result, want: %+v != got: %+v", got, items)
	}
}
