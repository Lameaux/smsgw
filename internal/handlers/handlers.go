package handlers

import (
	"strconv"
	"time"

	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/utils"
	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 10
)

func queryParamIntDefault(c *gin.Context, param string, def int) (int, error) {
	value := c.Query(param)
	if value == "" {
		return def, nil
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return def, err
	}

	return i, nil
}

func queryParamTime(c *gin.Context, param string) (*time.Time, error) {
	value := c.Query(param)
	if value == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func commonSearchParams(c *gin.Context) (*inputs.SearchParams, error) {
	offset, err := queryParamIntDefault(c, "offset", 0)
	if err != nil {
		return nil, err
	}

	limit, err := queryParamIntDefault(c, "limit", defaultLimit)
	if err != nil {
		return nil, err
	}

	createdAtFrom, err := queryParamTime(c, "created_at_from")
	if err != nil {
		return nil, err
	}

	createdAtTo, err := queryParamTime(c, "created_at_to")
	if err != nil {
		return nil, err
	}

	p := &inputs.SearchParams{
		Offset:        offset,
		Limit:         limit,
		CreatedAtFrom: createdAtFrom,
		CreatedAtTo:   createdAtTo,
	}

	return p, nil
}

func messageSearchParams(c *gin.Context) (*inputs.MessageParams, error) {
	p := &inputs.MessageParams{}

	msisdn := c.Query("msisdn")
	if msisdn != "" {
		msisdn, err := utils.NormalizeMSISDN(msisdn)
		if err != nil {
			return nil, err
		}

		p.MSISDN = &msisdn
	}

	status := c.Query("status")
	if status != "" {
		p.Status = &status
	}

	return p, nil
}
