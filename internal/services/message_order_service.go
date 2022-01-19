package services

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/utils"
	"euromoby.com/smsgw/internal/views"
)

type MessageOrderService struct {
	app *config.AppConfig
}

func NewMessageOrderService(app *config.AppConfig) *MessageOrderService {
	return &MessageOrderService{app}
}

func (s *MessageOrderService) FindByMerchantAndID(merchantID, id string) (*views.MessageOrderStatus, error) {
	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := s.app.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	messageOrderRepo := repos.NewMessageOrderRepo(conn)

	messageOrder, err := messageOrderRepo.FindByMerchantAndID(merchantID, id)
	if err != nil {
		return nil, err
	}

	if messageOrder == nil {
		return nil, nil
	}

	outboundMessageRepo := repos.NewOutboundMessageRepo(conn)

	messages, err := outboundMessageRepo.FindByMessageOrderID(merchantID, id)
	if err != nil {
		return nil, err
	}

	return views.NewMessageOrderStatus(messageOrder, messages), nil
}

func (s *MessageOrderService) FindByQuery(p *inputs.MessageOrderSearchParams) ([]*models.MessageOrder, error) {
	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := s.app.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	messageOrderRepo := repos.NewMessageOrderRepo(conn)

	messageOrders, err := messageOrderRepo.FindByQuery(p)
	if err != nil {
		return nil, err
	}

	return messageOrders, nil
}

func (s *MessageOrderService) SendMessage(params *inputs.SendMessageParams) (*views.MessageOrderStatus, error) {
	ctx, cancel := repos.DBTxContext()
	defer cancel()

	tx, err := s.app.DBPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	order := s.makeMessageOrder(params)
	err = messageOrderRepo.Save(order)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	var messages []*models.OutboundMessage

	for _, msisdn := range params.To {
		outboundMessage := models.NewOutboundMessage(params.MerchantID, order.ID, msisdn)
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

func (s *MessageOrderService) makeMessageOrder(p *inputs.SendMessageParams) *models.MessageOrder {
	now := utils.Now()
	return &models.MessageOrder{
		MerchantID:          p.MerchantID,
		Sender:              p.Sender,
		Body:                p.Body,
		NotificationURL:     p.NotificationURL,
		ClientTransactionID: p.ClientTransactionID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
