package main

import "testing"

func TestNewDefaultOTP(t *testing.T) {
	for _, c := range []struct {
		secret     string
		counters   []uint64
		wantValues []string
	}{
		{
			secret:   "12345678901234567890",
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
	} {
		if len(c.counters) != len(c.wantValues) {
			panic("invalid test table: counters and values len mismatch")
		}

		otp := NewDefaultOTP(c.secret)

		for i, counter := range c.counters {
			if got := otp.Generate(counter); got != c.wantValues[i] {
				t.Errorf("wrong otp value: %s", got)
			}
		}
	}
}
