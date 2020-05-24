package main

import (
	"flag"
	"fmt"
)

var secret string
var at int64

func init() {
	flag.StringVar(&secret, "secret", "", "Your two-factor secret")
	flag.Int64Var(&at, "at", 0, "")
}

func main() {
	flag.Parse()

	secret = DecodeBase32Secret(secret)

	totp := NewDefaultTOTP(secret)
	fmt.Printf("secret: %s\n", secret)

	var token string

	if at == 0 {
		token = totp.Now()
	} else {
		token = totp.At(at)
	}

	fmt.Println(token)
}
