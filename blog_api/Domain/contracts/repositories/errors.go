package repositories

import "errors"

// common error in repository
var (
	ErrNotFound = errors.New("not found")
)