package utils

import "errors"

var (
	ErrEmailOrUsernameTaken = errors.New("email or username already taken")
)
