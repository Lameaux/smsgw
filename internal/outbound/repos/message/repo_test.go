package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"euromoby.com/smsgw/internal/outbound/models"

	coremodels "euromoby.com/core/models"
)

func TestRepo_Save(t *testing.T) {
	r := NewRepo(TestApp.Config.DBPool)

	// New outbound message
	om := models.NewMessage(coremodels.NewUUID(), coremodels.NewUUID(), 420123456789)
	err := r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	assert.NotNil(t, om.ID, "ID is not set")
	assert.Equal(t, models.MessageStatusNew, om.Status, "Invalid status")
}

func TestRepo_FindByMerchantAndID(t *testing.T) {
	r := NewRepo(TestApp.Config.DBPool)

	// Unknown IDs
	_, err := r.FindByMerchantAndID(coremodels.NewUUID(), coremodels.NewUUID())
	assert.Equal(t, coremodels.ErrNotFound, err, "ErrNotFound expected")

	// Populate database
	merchantID := coremodels.NewUUID()
	om := models.NewMessage(merchantID, coremodels.NewUUID(), 420123456789)
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	// Find message
	found, err := r.FindByMerchantAndID(merchantID, om.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "OutboundMessage not found")
	assert.Equal(t, om, found, "OutboundMessage is not equal")

	// Invalid merchant ID
	_, err = r.FindByMerchantAndID(coremodels.NewUUID(), om.ID)
	assert.Equal(t, coremodels.ErrNotFound, err, "ErrNotFound expected")
}

func TestRepo_FindByMerchantAndGroupID(t *testing.T) {
	r := NewRepo(TestApp.Config.DBPool)

	// Unknown IDs
	found, err := r.FindByMerchantAndGroupID(coremodels.NewUUID(), coremodels.NewUUID())
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.Len(t, found, 0, "result should be empty")

	// Populate database
	merchantID := coremodels.NewUUID()
	messageGroupID := coremodels.NewUUID()
	om := models.NewMessage(merchantID, messageGroupID, 420123456789)
	err = r.Save(om)
	require.NoError(t, err, "Saving OutboundMessage failed")

	// Find by merchant and messageGroupID
	found, err = r.FindByMerchantAndGroupID(merchantID, messageGroupID)
	require.NoError(t, err, "Finding OutboundMessage failed")
	assert.NotEmpty(t, found, "OutboundMessage not found")
	assert.Equal(t, om, found[0], "OutboundMessage is not equal")
}
