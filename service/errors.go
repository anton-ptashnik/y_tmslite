package service

import "errors"

var (
	errLastStatusDelAttempt = errors.New("last status cannot be deleted")
)