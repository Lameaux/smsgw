package models

import (
	"errors"
	"regexp"
	"strconv"
)

type MSISDN int64

var (
	msisdnRegex      = regexp.MustCompile(`^(\+|00)?([1-9]\d{7,14})$`)
	ErrInvalidMSISDN = errors.New("invalid msisdn format")
)

func NormalizeMSISDN(msisdn string) (MSISDN, error) {
	match := msisdnRegex.FindStringSubmatch(msisdn)
	if match == nil {
		return 0, ErrInvalidMSISDN
	}

	normalized, err := strconv.ParseInt(match[2], 10, 64)
	if err != nil {
		return 0, ErrInvalidMSISDN
	}

	return MSISDN(normalized), nil
}
