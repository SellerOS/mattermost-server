package api4

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
)

func TestGetUserStatus(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	userStatus, resp := Client.GetUserStatus(th.BasicUser.ClientId, "")
	CheckNoError(t, resp)
	if userStatus.Status != "offline" {
		t.Fatal("Should return offline status")
	}

	th.App.SetStatusOnline(th.BasicUser.ClientId, true)
	userStatus, resp = Client.GetUserStatus(th.BasicUser.ClientId, "")
	CheckNoError(t, resp)
	if userStatus.Status != "online" {
		t.Fatal("Should return online status")
	}

	th.App.SetStatusAwayIfNeeded(th.BasicUser.ClientId, true)
	userStatus, resp = Client.GetUserStatus(th.BasicUser.ClientId, "")
	CheckNoError(t, resp)
	if userStatus.Status != "away" {
		t.Fatal("Should return away status")
	}

	th.App.SetStatusDoNotDisturb(th.BasicUser.ClientId)
	userStatus, resp = Client.GetUserStatus(th.BasicUser.ClientId, "")
	CheckNoError(t, resp)
	if userStatus.Status != "dnd" {
		t.Fatal("Should return dnd status")
	}

	th.App.SetStatusOffline(th.BasicUser.ClientId, true)
	userStatus, resp = Client.GetUserStatus(th.BasicUser.ClientId, "")
	CheckNoError(t, resp)
	if userStatus.Status != "offline" {
		t.Fatal("Should return offline status")
	}

	//Get user2 status logged as user1
	userStatus, resp = Client.GetUserStatus(th.BasicUser2.ClientId, "")
	CheckNoError(t, resp)
	if userStatus.Status != "offline" {
		t.Fatal("Should return offline status")
	}

	Client.Logout()

	_, resp = Client.GetUserStatus(th.BasicUser2.ClientId, "")
	CheckUnauthorizedStatus(t, resp)

	th.LoginBasic2()
	userStatus, resp = Client.GetUserStatus(th.BasicUser2.ClientId, "")
	CheckNoError(t, resp)
	if userStatus.Status != "offline" {
		t.Fatal("Should return offline status")
	}
}

func TestGetUsersStatusesByIds(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	ClientIds := []string{th.BasicUser.ClientId, th.BasicUser2.ClientId}

	usersStatuses, resp := Client.GetUsersStatusesByIds(ClientIds)
	CheckNoError(t, resp)
	for _, userStatus := range usersStatuses {
		if userStatus.Status != "offline" {
			t.Fatal("Status should be offline")
		}
	}

	th.App.SetStatusOnline(th.BasicUser.ClientId, true)
	th.App.SetStatusOnline(th.BasicUser2.ClientId, true)
	usersStatuses, resp = Client.GetUsersStatusesByIds(ClientIds)
	CheckNoError(t, resp)
	for _, userStatus := range usersStatuses {
		if userStatus.Status != "online" {
			t.Fatal("Status should be offline")
		}
	}

	th.App.SetStatusAwayIfNeeded(th.BasicUser.ClientId, true)
	th.App.SetStatusAwayIfNeeded(th.BasicUser2.ClientId, true)
	usersStatuses, resp = Client.GetUsersStatusesByIds(ClientIds)
	CheckNoError(t, resp)
	for _, userStatus := range usersStatuses {
		if userStatus.Status != "away" {
			t.Fatal("Status should be offline")
		}
	}

	th.App.SetStatusDoNotDisturb(th.BasicUser.ClientId)
	th.App.SetStatusDoNotDisturb(th.BasicUser2.ClientId)
	usersStatuses, resp = Client.GetUsersStatusesByIds(ClientIds)
	CheckNoError(t, resp)
	for _, userStatus := range usersStatuses {
		if userStatus.Status != "dnd" {
			t.Fatal("Status should be offline")
		}
	}

	Client.Logout()

	_, resp = Client.GetUsersStatusesByIds(ClientIds)
	CheckUnauthorizedStatus(t, resp)
}

func TestUpdateUserStatus(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	toUpdateUserStatus := &model.Status{Status: "online", UserId: th.BasicUser.ClientId}
	updateUserStatus, resp := Client.UpdateUserStatus(th.BasicUser.ClientId, toUpdateUserStatus)
	CheckNoError(t, resp)
	if updateUserStatus.Status != "online" {
		t.Fatal("Should return online status")
	}

	toUpdateUserStatus.Status = "away"
	updateUserStatus, resp = Client.UpdateUserStatus(th.BasicUser.ClientId, toUpdateUserStatus)
	CheckNoError(t, resp)
	if updateUserStatus.Status != "away" {
		t.Fatal("Should return away status")
	}

	toUpdateUserStatus.Status = "dnd"
	updateUserStatus, resp = Client.UpdateUserStatus(th.BasicUser.ClientId, toUpdateUserStatus)
	CheckNoError(t, resp)
	if updateUserStatus.Status != "dnd" {
		t.Fatal("Should return dnd status")
	}

	toUpdateUserStatus.Status = "offline"
	updateUserStatus, resp = Client.UpdateUserStatus(th.BasicUser.ClientId, toUpdateUserStatus)
	CheckNoError(t, resp)
	if updateUserStatus.Status != "offline" {
		t.Fatal("Should return offline status")
	}

	toUpdateUserStatus.Status = "online"
	toUpdateUserStatus.UserId = th.BasicUser2.ClientId
	_, resp = Client.UpdateUserStatus(th.BasicUser2.ClientId, toUpdateUserStatus)
	CheckForbiddenStatus(t, resp)

	toUpdateUserStatus.Status = "online"
	updateUserStatus, _ = th.SystemAdminClient.UpdateUserStatus(th.BasicUser2.ClientId, toUpdateUserStatus)
	if updateUserStatus.Status != "online" {
		t.Fatal("Should return online status")
	}

	_, resp = Client.UpdateUserStatus(th.BasicUser.ClientId, toUpdateUserStatus)
	CheckBadRequestStatus(t, resp)

	Client.Logout()

	_, resp = Client.UpdateUserStatus(th.BasicUser2.ClientId, toUpdateUserStatus)
	CheckUnauthorizedStatus(t, resp)
}
