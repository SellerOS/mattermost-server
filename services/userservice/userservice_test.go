package userservice

import (
	"testing"
	"github.com/mattermost/mattermost-server/model"
)

func TestSaveUser(t *testing.T) {
	t.Run("SaveUser", func(t *testing.T) {
		user := model.UserIms{ClientId: "18071711142411212", FirstName: "test", LastName: "wengang"}
		SaveUser(&user)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("UpdateUser", func(t *testing.T) {
		user := model.UserIms{ClientId: "18071711142411212", FirstName: "test222", LastName: "wengang222"}
		UpdateUser(&user)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("GetUser", func(t *testing.T) {
		GetUser("180717111424112312")
	})
}
