package totp

import (
	"testing"
)

func TestDecodeBase32Secret(t *testing.T) {
	for _, c := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
		{
			name:  "test padding #1",
			input: "GE",
			want:  "1",
		},
		{
			name:  "test padding #2",
			input: "GE=====",
			want:  "1",
		},
		{
			name:  "test lowercase",
			input: "ge=====",
			want:  "1",
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := DecodeBase32Secret(c.input); got != c.want {
				t.Errorf("wrong decoded value for %s input, want: %s != got: %s", c.input, c.want, got)
			}
		})
	}
}
