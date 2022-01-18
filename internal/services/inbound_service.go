package services

import (
	"euromoby.com/smsgw/internal/config"
	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/repos"
	"euromoby.com/smsgw/internal/utils"
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
	return inboundMessageRepo.FindByID(shortcode, id)
}

func (s *InboundService) AckByShortcodeAndID(shortcode, id string) (*models.InboundMessage, error) {
	ctx, cancel := repos.DBTxContext()
	defer cancel()

	tx, err := s.app.DBPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	inboundMessageRepo := repos.NewInboundMessageRepo(tx)
	message, err := inboundMessageRepo.FindByID(shortcode, id)

	if err != nil || message == nil {
		tx.Rollback(ctx)
		return message, err
	}

	if message.Status == models.InboundMessageStatusDelivered {
		tx.Rollback(ctx)
		return nil, models.ErrAlreadyAcked
	}

	message.Status = models.InboundMessageStatusDelivered
	message.UpdatedAt = utils.Now()
	err = inboundMessageRepo.UpdateStatus(message)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	notificationRepo := repos.NewDeliveryNotificationRepo(tx)

	n := models.MakeInboundDeliveryNotification(message)
	err = notificationRepo.Save(n)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	tx.Commit(ctx)

	return message, err
}
