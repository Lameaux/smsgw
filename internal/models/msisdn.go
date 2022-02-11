package models

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type MSISDN int64

var (
	msisdnRegex      = regexp.MustCompile(`^(\+|00)?([1-9]\d{7,14})$`)
	ErrInvalidMSISDN = errors.New("invalid msisdn format")
)

func (msisdn MSISDN) String() string {
	return strconv.FormatInt(int64(msisdn), 10) //nolint:gomnd
}

func stringToMSISDN(s string) (MSISDN, error) {
	if len(s) < 8 || len(s) > 15 {
		return 0, ErrInvalidMSISDN
	}

	msisdn, err := strconv.ParseInt(s, 10, 64) //nolint:gomnd
	if err != nil {
		return 0, ErrInvalidMSISDN
	}

	return MSISDN(msisdn), nil
}

func NormalizeMSISDNRegex(msisdn string) (MSISDN, error) {
	match := msisdnRegex.FindStringSubmatch(msisdn)
	if match == nil {
		return 0, ErrInvalidMSISDN
	}

	return stringToMSISDN(match[2])
}

func NormalizeMSISDN(msisdn string) (MSISDN, error) {
	switch {
	case strings.HasPrefix(msisdn, "00+") || strings.HasPrefix(msisdn, "+00"):
		return 0, ErrInvalidMSISDN
	case strings.HasPrefix(msisdn, "+"):
		return stringToMSISDN(msisdn[1:])
	case strings.HasPrefix(msisdn, "00"):
		return stringToMSISDN(msisdn[2:])
	default:
		return stringToMSISDN(msisdn)
	}
}
