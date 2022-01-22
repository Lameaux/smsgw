package handlers

import "errors"

var (
	ErrMessageNotFound      = errors.New("message not found")
	ErrMessageOrderNotFound = errors.New("message order not found")
)
