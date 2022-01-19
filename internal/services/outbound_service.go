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

	message, err := outboundMessageRepo.FindByID(merchantID, id)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, nil
	}

	messageOrderRepo := repos.NewMessageOrderRepo(conn)

	messageOrder, err := messageOrderRepo.FindByMerchantAndID(merchantID, message.MessageOrderID)
	if err != nil {
		return nil, err
	}

	if messageOrder == nil {
		return nil, nil
	}

	return views.NewOutboundMessageDetail(message, messageOrder), nil
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
