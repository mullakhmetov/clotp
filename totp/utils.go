package totp

import (
	"encoding/base32"
	"strings"
)

func DecodeBase32Secret(s string) string {
	s = strings.ToUpper(s)

	missingPadding := len(s) % 8
	if missingPadding != 0 {
		s = s + strings.Repeat("=", 8-missingPadding)
	}
	bytes, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		panic("decode secret failed")
	}
	return string(bytes)
}
