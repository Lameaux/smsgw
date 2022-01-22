package services

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/inputs"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
)

type InboundService struct {
	app *config.AppConfig
}

func NewInboundService(app *config.AppConfig) *InboundService {
	return &InboundService{app}
}

func (s *InboundService) ValidateShortcode(merchantID, shortcode string) error {
	return nil
}

func (s *InboundService) FindByShortcodeAndID(shortcode, id string) (*models.InboundMessage, error) {
	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := s.app.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	inboundMessageRepo := repos.NewInboundMessageRepo(conn)
	return inboundMessageRepo.FindByShortcodeAndID(shortcode, id)
}

func (s *InboundService) SaveMessage(m *models.InboundMessage) error {
	ctx, done := repos.DBConnContext()
	defer done()

	conn, err := s.app.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	inboundMessageRepo := repos.NewInboundMessageRepo(conn)
	return inboundMessageRepo.Save(m)
}

func (s *InboundService) AckByShortcodeAndID(shortcode, id string) (*models.InboundMessage, error) {
	ctx, cancel := repos.DBTxContext()
	defer cancel()

	tx, err := s.app.DBPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	inboundMessageRepo := repos.NewInboundMessageRepo(tx)
	message, err := inboundMessageRepo.FindByShortcodeAndID(shortcode, id)

	if err != nil || message == nil {
		return message, err
	}

	if message.Status == models.InboundMessageStatusDelivered {
		return nil, models.ErrAlreadyAcked
	}

	message.Status = models.InboundMessageStatusDelivered
	err = inboundMessageRepo.Update(message)
	if err != nil {
		return nil, err
	}

	notificationRepo := repos.NewDeliveryNotificationRepo(tx)

	n := models.MakeInboundDeliveryNotification(message)
	err = notificationRepo.Save(n)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return message, err
}

func (s *InboundService) FindByQuery(p *inputs.InboundMessageSearchParams) ([]*models.InboundMessage, error) {
	ctx, cancel := repos.DBConnContext()
	defer cancel()

	conn, err := s.app.DBPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	inboundMessageRepo := repos.NewInboundMessageRepo(conn)

	messages, err := inboundMessageRepo.FindByQuery(p)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
