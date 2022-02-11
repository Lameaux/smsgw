package services

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/views"
)

type OutboundService struct {
	app *config.AppConfig
}

func NewOutboundService(app *config.AppConfig) *OutboundService {
	return &OutboundService{app}
}

func (s *OutboundService) FindByMerchantAndID(merchantID, id string) (*views.OutboundMessageDetail, error) {
	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := s.app.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	outboundMessageRepo := repos.NewOutboundMessageRepo(conn)

	message, err := outboundMessageRepo.FindByMerchantAndID(merchantID, id)
	if err != nil {
		return nil, err
	}

	messageOrderRepo := repos.NewMessageOrderRepo(conn)

	messageOrder, err := messageOrderRepo.FindByID(message.MessageOrderID)
	if err != nil {
		return nil, err
	}

	return views.NewOutboundMessageDetail(message, messageOrder), nil
}

func (s *OutboundService) AckByProviderAndMessageID(providerID, messageID string) (*models.OutboundMessage, error) {
	ctx, cancel := repos.DBTxContext()
	defer cancel()

	tx, err := s.app.DBPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	outboundMessageRepo := repos.NewOutboundMessageRepo(tx)

	message, err := outboundMessageRepo.FindByProviderAndMessageID(providerID, messageID)
	if err != nil {
		return nil, err
	}

	messageOrderRepo := repos.NewMessageOrderRepo(tx)

	messageOrder, err := messageOrderRepo.FindByID(message.MessageOrderID)
	if err != nil {
		return nil, err
	}

	if message.Status == models.OutboundMessageStatusDelivered {
		return nil, models.ErrAlreadyAcked
	}

	message.Status = models.OutboundMessageStatusDelivered

	err = outboundMessageRepo.Update(message)
	if err != nil {
		return nil, err
	}

	if messageOrder.NotificationURL != nil {
		notificationRepo := repos.NewDeliveryNotificationRepo(tx)
		n := models.MakeOutboundDeliveryNotification(message)

		err = notificationRepo.Save(n)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return message, err
}

func (s *OutboundService) FindByQuery(p *inputs.OutboundMessageSearchParams) ([]*models.OutboundMessage, error) {
	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := s.app.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	outboundMessageRepo := repos.NewOutboundMessageRepo(conn)

	messages, err := outboundMessageRepo.FindByQuery(p)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
