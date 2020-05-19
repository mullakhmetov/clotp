package main

import "testing"

func TestNewDefaultTOTP(t *testing.T) {
	for _, c := range []struct {
		secret     string
		timestamps []int64
		wantValues []string
	}{
		{
			secret:     "12345678901234567890",
			timestamps: []int64{59},
			wantValues: []string{"287082"},
		},
	} {
		if len(c.timestamps) != len(c.wantValues) {
			panic("invalid test table: counters and values len mismatch")
		}

		otp := NewDefaultTOTP(c.secret)

		for i, ts := range c.timestamps {
			if got := otp.At(ts); got != c.wantValues[i] {
				t.Errorf("wrong otp value: %s", got)
			}
		}
	}
}
