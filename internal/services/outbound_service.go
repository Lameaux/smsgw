package services

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/utils"
	"euromoby.com/smsgw/internal/views"
)

type OutboundService struct {
	app *config.AppConfig
}

func NewOutboundService(app *config.AppConfig) *OutboundService {
	return &OutboundService{app}
}

func (s *OutboundService) SendMessage(merchantID string, params *views.SendMessageParams) (*views.MessageOrderStatus, error) {
	ctx, cancel := repos.DBTxContext()
	defer cancel()

	tx, err := s.app.DBPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	order := s.makeMessageOrder(merchantID, params)
	err = messageOrderRepo.Save(order)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	var messages []*models.OutboundMessage

	for _, msisdn := range params.To {
		outboundMessage := models.NewOutboundMessage(merchantID, order.ID, msisdn)
		messages = append(messages, outboundMessage)
	}

	outboundMessageRepo := repos.NewOutboundMessageRepo(tx)

	for i := range messages {
		err = outboundMessageRepo.Save(messages[i])
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
	}

	tx.Commit(ctx)

	return views.NewMessageOrderStatus(order, messages), nil
}

func (s *OutboundService) makeMessageOrder(merchantID string, p *views.SendMessageParams) *models.MessageOrder {
	now := utils.Now()
	return &models.MessageOrder{
		MerchantID:          merchantID,
		Sender:              p.Sender,
		Body:                p.Body,
		NotificationURL:     p.NotificationURL,
		ClientTransactionID: p.ClientTransactionID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
