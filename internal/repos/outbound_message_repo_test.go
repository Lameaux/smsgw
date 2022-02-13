package repos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"euromoby.com/smsgw/internal/models"
)

func TestOutboundMessageRepo_Save(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	// New outbound message
	om := models.NewOutboundMessage(models.NewUUID(), models.NewUUID(), 420123456789)
	err := r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	assert.NotNil(t, om.ID, "ID is not set")
	assert.Equal(t, models.OutboundMessageStatusNew, om.Status, "Invalid status")
}

func TestOutboundMessageRepo_FindByMerchantAndID(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	// Unknown IDs
	_, err := r.FindByMerchantAndID(models.NewUUID(), models.NewUUID())
	assert.Equal(t, models.ErrNotFound, err, "ErrNotFound expected")

	// Populate database
	merchantID := models.NewUUID()
	om := models.NewOutboundMessage(merchantID, models.NewUUID(), 420123456789)
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	// Find message
	found, err := r.FindByMerchantAndID(merchantID, om.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "OutboundMessage not found")
	assert.Equal(t, om, found, "OutboundMessage is not equal")

	// Invalid merchant ID
	_, err = r.FindByMerchantAndID(models.NewUUID(), om.ID)
	assert.Equal(t, models.ErrNotFound, err, "ErrNotFound expected")
}

func TestOutboundMessageRepo_FindByMerchantAndOrderID(t *testing.T) {
	r := NewOutboundMessageRepo(TestAppConfig.DBPool)

	// Unknown IDs
	found, err := r.FindByMerchantAndOrderID(models.NewUUID(), models.NewUUID())
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.Len(t, found, 0, "result should be empty")

	// Populate database
	merchantID := models.NewUUID()
	messageOrderID := models.NewUUID()
	om := models.NewOutboundMessage(merchantID, messageOrderID, 420123456789)
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	// Find by merchant and messageOrderID
	found, err = r.FindByMerchantAndOrderID(merchantID, messageOrderID)
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.NotEmpty(t, found, "OutboundMessage not found")
	assert.Equal(t, om, found[0], "OutboundMessage is not equal")
}
