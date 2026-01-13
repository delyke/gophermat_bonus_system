package model

import "errors"

var (
	ErrUserNotFound   = errors.New("not found")
	ErrBadCredentials = errors.New("bad credentials")
	ErrLoginTaken     = errors.New("login already taken")
)
