package config

import (
	"bytes"
	"reflect"
	"testing"
)

var (
	fullConfig = []byte(`
	[Name-1]
	issuer=issuer-1
	secret=secret-key-1
	algorithm=sha1
	digits=6
	`)

	onlySecret = []byte(`
	[Name-2]
	secret=secret-key-2
	`)

	noSecret = []byte(`
	[Name-3]
	issuer=issuer-3
	algorithm=sha1
	digits=6
	`)

	multiple = []byte(`
	[Name-4]
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

func TestLoadConfig(t *testing.T) {
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
			got, err := loadConfig(bytes.NewReader(c.input))
			if err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("wrong config loading, want: %+v != got: %+v", c.want, got)
			}
		})
	}
}
