package userErrors

import "errors"

var (
	ErrMissingUserID = errors.New("missing clerk user id")
	ErrUserNotFound  = errors.New("user not found")
)
