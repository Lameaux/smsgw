package handlers

import (
	commoninputs "github.com/Lameaux/smsgw/internal/common/inputs"
	"time"

	"github.com/gin-gonic/gin"

	coremodels "github.com/Lameaux/core/models"
	"github.com/Lameaux/core/utils"
)

const (
	defaultLimit = 10
)

func QueryParamUint64Default(c *gin.Context, param string, def uint64) (uint64, error) {
	value := c.Query(param)
	if value == "" {
		return def, nil
	}

	i, err := utils.ParseUint64(value)
	if err != nil {
		return def, err
	}

	return i, nil
}

func QueryParamTime(c *gin.Context, param string) (*time.Time, error) {
	value := c.Query(param)
	if value == "" {
		return nil, nil //nolint: nilnil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func CommonSearchParams(c *gin.Context) (*commoninputs.SearchParams, error) {
	offset, err := QueryParamUint64Default(c, "offset", 0)
	if err != nil {
		return nil, err
	}

	limit, err := QueryParamUint64Default(c, "limit", defaultLimit)
	if err != nil {
		return nil, err
	}

	createdAtFrom, err := QueryParamTime(c, "created_at_from")
	if err != nil {
		return nil, err
	}

	createdAtTo, err := QueryParamTime(c, "created_at_to")
	if err != nil {
		return nil, err
	}

	p := commoninputs.SearchParams{
		Offset:        offset,
		Limit:         limit,
		CreatedAtFrom: createdAtFrom,
		CreatedAtTo:   createdAtTo,
	}

	return &p, nil
}

func MessageSearchParams(c *gin.Context) (*commoninputs.MessageParams, error) {
	p := commoninputs.MessageParams{}

	if msisdn := c.Query("msisdn"); msisdn != "" {
		normalized, err := coremodels.NormalizeMSISDN(msisdn)
		if err != nil {
			return nil, err
		}

		p.MSISDN = &normalized
	}

	if status := c.Query("status"); status != "" {
		p.Status = &status
	}

	return &p, nil
}
