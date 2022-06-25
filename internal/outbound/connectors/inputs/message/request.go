package message

import coremodels "github.com/Lameaux/core/models"

type Request struct {
	MSISDN              coremodels.MSISDN
	Sender              string
	Body                string
	ClientTransactionID string
}
