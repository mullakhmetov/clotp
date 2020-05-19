package main

import (
	"time"
)

const defaultTimeStep uint64 = 30

// DefaultOTP returns 6-digits HMAC-SHA-1 OTP based on given secret and counter
func NewDefaultTOTP(secret string) *TOTP {
	return &TOTP{NewDefaultOTP(secret), defaultTimeStep}
}

// NewOTP returns OTP object
func NewTOTP(opts Opts, timestep uint64) *TOTP {
	return &TOTP{NewOTP(opts), timestep}
}

type TOTP struct {
	*OTP
	TimeStep uint64
}

func (t *TOTP) Now() string {
	ts := time.Now().Unix()
	return t.At(ts)
}

func (t *TOTP) At(ts int64) string {
	return t.Generate(uint64(ts / int64(t.TimeStep)))
}
