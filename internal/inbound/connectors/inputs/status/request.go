package status

import coremodels "github.com/Lameaux/core/models"

type Request struct {
	Provider            string
	MessageID           string
	MSISDN              coremodels.MSISDN
	Status              string
	ClientTransactionID string
}
