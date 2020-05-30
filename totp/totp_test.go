package totp

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"testing"
	"time"
)

func TestTOTP_At(t *testing.T) {
	for _, c := range []struct {
		name       string
		opts       Opts
		step       int
		timestamps []int64
		wantValues []string
	}{
		{
			name: "sha1",
			opts: Opts{
				Digits:    8,
				Secret:    "12345678901234567890",
				Algorithm: sha1.New,
			},
			step:       30,
			timestamps: []int64{59, 1111111109, 1111111111, 1234567890, 2000000000, 20000000000},
			wantValues: []string{"94287082", "07081804", "14050471", "89005924", "69279037", "65353130"},
		},
		{
			name: "sha1",
			opts: Opts{
				Digits:    8,
				Secret:    "12345678901234567890",
				Algorithm: sha1.New,
			},
			step:       1,
			timestamps: []int64{59, 1111111109, 1111111111, 1234567890, 2000000000, 20000000000},
			wantValues: []string{"24083773", "14510233", "54073950", "03965462", "21054266", "04468884"},
		},
		{
			name: "sha256",
			opts: Opts{
				Digits:    6,
				Secret:    "12345678901234567890",
				Algorithm: sha256.New,
			},
			step:       30,
			timestamps: []int64{59, 1111111109, 1111111111, 1234567890, 2000000000, 20000000000},
			wantValues: []string{"247374", "756375", "584430", "829826", "428693", "142410"},
		},
		{
			name: "sha512",
			opts: Opts{
				Digits:    8,
				Secret:    "12345678901234567890",
				Algorithm: sha512.New,
			},
			step:       30,
			timestamps: []int64{59, 1111111109, 1111111111, 1234567890, 2000000000, 20000000000},
			wantValues: []string{"69342147", "63049338", "54380122", "76671578", "56464532", "69481994"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if len(c.timestamps) != len(c.wantValues) {
				panic("invalid test table: counters and values len mismatch")
			}

			totp := NewTOTP(c.opts, c.step)

			for i, ts := range c.timestamps {
				if got := totp.At(ts); got != c.wantValues[i] {
					t.Errorf("wrong otp value for %d counter, want: %s != got: %s", ts, c.wantValues[i], got)
				}
			}
		})
	}
}

func TestTOTP_Now(t *testing.T) {
	totp := NewTOTP(Opts{Digits: 6, Secret: "12345678901234567890", Algorithm: sha1.New}, 30)

	if totp.Now() != totp.At(time.Now().Unix()) {
		t.Errorf("wrong Now() value")
	}
}

func TestAdfs(t *testing.T) {
	totp := NewTOTP(Opts{Digits: 6, Secret: "tfci3jrtwuwdjlyo", Algorithm: sha1.New}, 30)
	got := totp.Now()
	panic(got)
}
