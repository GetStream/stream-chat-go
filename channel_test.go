package stream_chat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateChannel(t *testing.T) {
	c := initClient(t)

	t.Run("get existing channel", func(t *testing.T) {
		ch := initChannel(t, c)
		got, err := c.CreateChannel(ch.Type, ch.ID, "gandalf", nil)
		mustNoError(t, err)

		assert.Equal(t, c, got.client)
		assert.Equal(t, ch.Type, got.Type)
		assert.Equal(t, ch.ID, got.ID)
		assert.Equal(t, got.MemberCount, ch.MemberCount)
		assert.Len(t, got.Members, got.MemberCount)
	})

	tests := []struct {
		_type   string
		id      string
		userID  string
		data    map[string]interface{}
		wantErr bool
	}{
		{"messaging", "mates", "gandalf", nil, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("create new channel %s:%s", tt._type, tt.id), func(t *testing.T) {
			got, err := c.CreateChannel(tt._type, tt.id, tt.userID, tt.data)
			if tt.wantErr {
				mustError(t, err)
			} else {
				mustNoError(t, err)
			}

			assert.Equal(t, tt._type, got.Type)
			assert.Equal(t, tt.id, got.ID)
			assert.Equal(t, tt.userID, got.CreatedBy.ID)
		})
	}
}

func TestChannel_AddMembers(t *testing.T) {
}

func TestChannel_AddModerators(t *testing.T) {

}

func TestChannel_BanUser(t *testing.T) {

}

func TestChannel_Create(t *testing.T) {

}

func TestChannel_Delete(t *testing.T) {

}

func TestChannel_DeleteReaction(t *testing.T) {

}

func TestChannel_DemoteModerators(t *testing.T) {

}

func TestChannel_GetReactions(t *testing.T) {

}

func TestChannel_GetReplies(t *testing.T) {

}

func TestChannel_MarkRead(t *testing.T) {

}

func TestChannel_Query(t *testing.T) {

}

func TestChannel_RemoveMembers(t *testing.T) {

}

func TestChannel_SendEvent(t *testing.T) {

}

func TestChannel_SendMessage(t *testing.T) {

}

func TestChannel_SendReaction(t *testing.T) {

}

func TestChannel_Truncate(t *testing.T) {

}

func TestChannel_UnBanUser(t *testing.T) {

}

func TestChannel_Update(t *testing.T) {

}

func Test_addUserID(t *testing.T) {
	id := "someid"
	params := map[string]interface{}{"test": 1}

	addUserID(params, id)

	assert.Contains(t, id, params["user"])
}
