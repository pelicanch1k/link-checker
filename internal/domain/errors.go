package domain

import "errors"

var (
	ErrEmptyURLs = errors.New("URLs list is empty")
	ErrTaskNotFound  = errors.New("task not found")
)