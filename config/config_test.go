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
	secret=secret-key
	algorithm=sha1
	digits=6
	`)
)

func TestLoadConfig(t *testing.T) {
	for _, c := range []struct {
		name  string
		input []byte
		want  Config
	}{
		{
			name:  "1",
			input: fullConfig,
			want: Config{Secrets: []secret{
				{Name: "Name-1", Issuer: "issuer-1", Key: "secret-key", Algorithm: "sha1", Digits: 6}},
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got, err := LoadConfig(bytes.NewReader(c.input))
			if err != nil {
				panic(err)
			}

			if !reflect.DeepEqual(c.want, *got) {
				t.Errorf("wrong config loading, want: %+v != got: %+v", c.want, got)
			}
		})
	}
}
