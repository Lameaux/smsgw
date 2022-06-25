package models

import "errors"

var (
	ErrSendFailed = errors.New("failed to send")
	ErrDeadLetter = errors.New("message can not be delivered")
)
