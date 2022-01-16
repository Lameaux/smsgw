package repos

import (
	"testing"

	"euromoby.com/smsgw/internal/models"
	"euromoby.com/smsgw/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageOrderRepo_Save(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)

	merchantID := utils.NewUUID()

	mo := models.NewSimpleMessageOrder(merchantID, "hello world")
	err := r.Save(mo)
	require.NoError(t, err, "Saving MessageOrder failed")

	assert.NotNil(t, mo.ID, "ID is not set")
	assert.NotNil(t, mo.ClientTransactionID, "ClientTransactionID ID is not set")
}

func TestMessageOrderRepo_FindByID(t *testing.T) {
	r := NewMessageOrderRepo(TestAppConfig.DBPool)

	found, err := r.FindByID(utils.NewUUID(), utils.NewUUID())
	assert.Nil(t, err, "unexpected error")
	assert.Nil(t, found, "wrong MessageOrder found")

	merchantID := utils.NewUUID()
	mo := models.NewSimpleMessageOrder(merchantID, "hello world")
	err = r.Save(mo)
	require.NoError(t, err, "Saving MessageOrder failed")

	found, err = r.FindByID(merchantID, mo.ID)
	assert.Nil(t, err, "unexpected error")
	assert.NotNil(t, found, "MessageOrder not found")

	assert.Equal(t, mo, found, "MessageOrder is not equal")
}
