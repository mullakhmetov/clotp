package main

import "errors"

var (
	ErrInvalidItem       = errors.New("item validation failed")
	ErrItemAlreadyExists = errors.New("item already exists")
)
