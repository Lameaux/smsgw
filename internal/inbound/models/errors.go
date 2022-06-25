package models

import "errors"

var (
	ErrMessageNotFound            = errors.New("message not found")
	ErrAlreadyAcked               = errors.New("message is already acked")
	ErrDuplicateProviderMessageID = errors.New("message_id already exists")
)
