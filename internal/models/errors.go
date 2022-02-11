package models

import "errors"

var (
	ErrAlreadyAcked                 = errors.New("message is already acked")
	ErrSendFailed                   = errors.New("failed to send")
	ErrInvalidJSON                  = errors.New("invalid json")
	ErrDeadLetter                   = errors.New("message can not be delivered")
	ErrDuplicateProviderMessageID   = errors.New("message_id already exists")
	ErrDuplicateClientTransactionID = errors.New("client_transaction_id already exists")
	ErrNotFound                     = errors.New("not found")
)
