package totp

import (
	"crypto/sha1" //nolint:gosec // used in hmac only, see RFC 4226 B.2. section
	"crypto/sha256"
	"crypto/sha512"
	"testing"
)

func TestGenerate(t *testing.T) {
	for _, c := range []struct {
		name       string
		opts       Opts
		counters   []uint64
		wantValues []string
	}{
		{
			name: "sha1 3 digits",
			opts: Opts{
				Digits:    3,
				Algorithm: sha1.New,
				Secret:    "12345678901234567890",
			},
			counters: []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantValues: []string{
				"224",
				"082",
				"152",
				"429",
				"314",
				"676",
				"922",
				"583",
				"871",
				"489",
			},
		},
		{
			name: "sha1 6 digits",
			opts: Opts{
				Digits:    6,
				Algorithm: sha1.New,
				Secret:    "12345678901234567890",
			},
			counters: []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantValues: []string{
				"755224",
				"287082",
				"359152",
				"969429",
				"338314",
				"254676",
				"287922",
				"162583",
				"399871",
				"520489",
			},
		},
		{
			name: "sha256 8 digits",
			opts: Opts{
				Digits:    8,
				Algorithm: sha256.New,
				Secret:    "12345678901234567890",
			},
			counters: []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantValues: []string{
				"74875740",
				"32247374",
				"66254785",
				"67496144",
				"25480556",
				"89697997",
				"40191609",
				"67579288",
				"83895912",
				"23184989",
			},
		},
		{
			name: "sha512 8 digits",
			opts: Opts{
				Digits:    8,
				Algorithm: sha512.New,
				Secret:    "12345678901234567890",
			},
			counters: []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantValues: []string{
				"04125165",
				"69342147",
				"71730102",
				"73778726",
				"81937510",
				"16848329",
				"36266680",
				"22588359",
				"45039399",
				"33643409",
			},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if len(c.counters) != len(c.wantValues) {
				panic("invalid test table: counters and values len mismatch")
			}

			otp := NewOTP(c.opts)

			for i, counter := range c.counters {
				if got := otp.Generate(counter); got != c.wantValues[i] {
					t.Errorf("wrong otp value, want: %s != got: %s", c.wantValues[i], got)
				}
			}
		})
	}
}

func TestNewOTP(t *testing.T) {
	opts := Opts{Digits: 1, Secret: "1", Algorithm: sha1.New}
	otp := NewOTP(opts)

	if otp.Digits != opts.Digits {
		t.Errorf("wrong digits: %d", otp.Digits)
	}

	if otp.Secret != opts.Secret {
		t.Errorf("wrong secret: %s", otp.Secret)
	}
}

func TestNewDefaultOTP(t *testing.T) {
	secret := "12345678901234567890"
	otp := NewDefaultOTP(secret)

	if otp.Digits != defaultDigits {
		t.Errorf("wrong digits: %d", otp.Digits)
	}

	if otp.Secret != secret {
		t.Errorf("wrong secret: %s", otp.Secret)
	}

	want := "755224"
	got := otp.Generate(0)

	if want != got {
		// seems like it's not the sha1 algorithm
		t.Errorf("wrong otp value for default otp: %s", got)
	}
}
