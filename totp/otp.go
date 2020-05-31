package totp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1" //nolint:gosec // used in hmac only, see RFC 4226 B.2. section
	"encoding/binary"
	"fmt"
	"hash"
	"math"
)

const defaultDigits = 6

type Opts struct {
	Digits    int
	Secret    string
	Algorithm func() hash.Hash
}

// DefaultOTP returns 6-digits HMAC-SHA-1 OTP based on given secret and counter
func NewDefaultOTP(secret string) *OTP {
	return &OTP{
		Opts{Digits: defaultDigits, Secret: secret, Algorithm: sha1.New},
	}
}

// NewOTP returns OTP object
func NewOTP(opts Opts) *OTP {
	return &OTP{opts}
}

type OTP struct {
	Opts
}

// Generate returns otp value based on given counter
func (o *OTP) Generate(counter uint64) string {
	mac := hmac.New(o.Algorithm, o.secret())
	if _, err := mac.Write(itob(counter)); err != nil {
		panic(err)
	}

	hmacResult := mac.Sum(nil)
	offset := int(hmacResult[len(hmacResult)-1] & 0xf)
	code := ((int(hmacResult[offset]) & 0x7f) << 24) |
		((int(hmacResult[offset+1] & 0xff)) << 16) |
		((int(hmacResult[offset+2] & 0xff)) << 8) |
		(int(hmacResult[offset+3]) & 0xff)

	value := code % int(math.Pow10(o.Digits))

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", o.Digits), value)
}

func (o *OTP) secret() []byte {
	return []byte(o.Secret)
}

func itob(i uint64) []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
