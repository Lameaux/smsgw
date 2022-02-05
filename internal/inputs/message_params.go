package inputs

import "euromoby.com/smsgw/internal/models"

type MessageParams struct {
	MSISDN *models.MSISDN
	Status *string
}
