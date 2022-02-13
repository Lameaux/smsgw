package repos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"euromoby.com/smsgw/internal/models"
)

func TestMessageOrderRepo_Save(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)
	merchantID := models.NewUUID()

	// New message order
	mo := models.NewSimpleMessageOrder(merchantID, "hello world")
	err := r.Save(mo)
	require.NoError(t, err, "Saving MessageOrder failed")
	assert.NotNil(t, mo.ID, "ID is not set")

	// Duplicate ClientTransactionID
	mo2 := models.NewSimpleMessageOrder(merchantID, "hello world")
	mo2.ClientTransactionID = mo.ClientTransactionID
	err = r.Save(mo2)
	require.ErrorAs(t, err, &models.ErrDuplicateClientTransactionID, "Duplicate MessageOrder created")
}

func TestMessageOrderRepo_FindByMerchantAndID(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)

	// Empty database
	_, err := r.FindByMerchantAndID(models.NewUUID(), models.NewUUID())
	assert.Equal(t, models.ErrNotFound, err, "ErrNotFound expected")

	// Populate database
	merchantID := models.NewUUID()
	mo := models.NewSimpleMessageOrder(merchantID, "hello world")
	err = r.Save(mo)
	require.NoError(t, err, "Saving MessageOrder failed")

	// Find order
	found, err := r.FindByMerchantAndID(merchantID, mo.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "MessageOrder not found")
	assert.Equal(t, mo, found, "MessageOrder is not equal")
}
