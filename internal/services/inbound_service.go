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
	tx, err := repos.Begin(s.app.DBPool)
	if err != nil {
		return nil, err
	}

	defer repos.Rollback(tx)

	inboundMessageRepo := repos.NewInboundMessageRepo(tx)

	message, err := inboundMessageRepo.FindByShortcodeAndID(shortcode, id)
	if err != nil {
		return nil, err
	}

	if message.Status == models.InboundMessageStatusDelivered {
		return nil, models.ErrAlreadyAcked
	}

	message.Status = models.InboundMessageStatusDelivered

	if err := inboundMessageRepo.Update(message); err != nil {
		return nil, err
	}

	notificationRepo := repos.NewDeliveryNotificationRepo(tx)

	n := models.MakeInboundDeliveryNotification(message)

	if err := notificationRepo.Save(n); err != nil {
		return nil, err
	}

	if err := repos.Commit(tx); err != nil {
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
