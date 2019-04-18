// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"strings"
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

func TestGetPreferences(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	th.LoginBasic()
	user1 := th.BasicUser

	category := model.NewId()
	preferences1 := model.Preferences{
		{
			UserId:   user1.ClientId,
			Category: category,
			Name:     model.NewId(),
		},
		{
			UserId:   user1.ClientId,
			Category: category,
			Name:     model.NewId(),
		},
		{
			UserId:   user1.ClientId,
			Category: model.NewId(),
			Name:     model.NewId(),
		},
	}

	Client.UpdatePreferences(user1.ClientId, &preferences1)

	prefs, resp := Client.GetPreferences(user1.ClientId)
	CheckNoError(t, resp)
	if len(prefs) != 4 {
		t.Fatal("received the wrong number of preferences")
	}

	for _, preference := range prefs {
		if preference.UserId != th.BasicUser.ClientId {
			t.Fatal("user id does not match")
		}
	}

	th.LoginBasic2()

	prefs, resp = Client.GetPreferences(th.BasicUser2.ClientId)
	CheckNoError(t, resp)

	if len(prefs) == 0 {
		t.Fatal("received the wrong number of preferences")
	}

	_, resp = Client.GetPreferences(th.BasicUser.ClientId)
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.GetPreferences(th.BasicUser2.ClientId)
	CheckUnauthorizedStatus(t, resp)
}

func TestGetPreferencesByCategory(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	th.LoginBasic()
	user1 := th.BasicUser

	category := model.NewId()
	preferences1 := model.Preferences{
		{
			UserId:   user1.ClientId,
			Category: category,
			Name:     model.NewId(),
		},
		{
			UserId:   user1.ClientId,
			Category: category,
			Name:     model.NewId(),
		},
		{
			UserId:   user1.ClientId,
			Category: model.NewId(),
			Name:     model.NewId(),
		},
	}

	Client.UpdatePreferences(user1.ClientId, &preferences1)

	prefs, resp := Client.GetPreferencesByCategory(user1.ClientId, category)
	CheckNoError(t, resp)

	if len(prefs) != 2 {
		t.Fatalf("received the wrong number of preferences %v:%v", len(prefs), 2)
	}

	_, resp = Client.GetPreferencesByCategory(user1.ClientId, "junk")
	CheckNotFoundStatus(t, resp)

	th.LoginBasic2()

	_, resp = Client.GetPreferencesByCategory(th.BasicUser2.ClientId, category)
	CheckNotFoundStatus(t, resp)

	_, resp = Client.GetPreferencesByCategory(user1.ClientId, category)
	CheckForbiddenStatus(t, resp)

	prefs, resp = Client.GetPreferencesByCategory(th.BasicUser2.ClientId, "junk")
	CheckNotFoundStatus(t, resp)

	if len(prefs) != 0 {
		t.Fatal("received the wrong number of preferences")
	}

	Client.Logout()
	_, resp = Client.GetPreferencesByCategory(th.BasicUser2.ClientId, category)
	CheckUnauthorizedStatus(t, resp)
}

func TestGetPreferenceByCategoryAndName(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	th.LoginBasic()
	user := th.BasicUser
	name := model.NewId()
	value := model.NewId()

	preferences := model.Preferences{
		{
			UserId:   user.ClientId,
			Category: model.PREFERENCE_CATEGORY_DIRECT_CHANNEL_SHOW,
			Name:     name,
			Value:    value,
		},
		{
			UserId:   user.ClientId,
			Category: model.PREFERENCE_CATEGORY_DIRECT_CHANNEL_SHOW,
			Name:     model.NewId(),
			Value:    model.NewId(),
		},
	}

	Client.UpdatePreferences(user.ClientId, &preferences)

	pref, resp := Client.GetPreferenceByCategoryAndName(user.ClientId, model.PREFERENCE_CATEGORY_DIRECT_CHANNEL_SHOW, name)
	CheckNoError(t, resp)

	if (pref.UserId != preferences[0].UserId) && (pref.Category != preferences[0].Category) && (pref.Name != preferences[0].Name) {
		t.Fatal("preference saved incorrectly")
	}

	preferences[0].Value = model.NewId()
	Client.UpdatePreferences(user.ClientId, &preferences)

	_, resp = Client.GetPreferenceByCategoryAndName(user.ClientId, "junk", preferences[0].Name)
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetPreferenceByCategoryAndName(user.ClientId, preferences[0].Category, "junk")
	CheckBadRequestStatus(t, resp)

	_, resp = Client.GetPreferenceByCategoryAndName(th.BasicUser2.ClientId, preferences[0].Category, "junk")
	CheckForbiddenStatus(t, resp)

	_, resp = Client.GetPreferenceByCategoryAndName(user.ClientId, preferences[0].Category, preferences[0].Name)
	CheckNoError(t, resp)

	Client.Logout()
	_, resp = Client.GetPreferenceByCategoryAndName(user.ClientId, preferences[0].Category, preferences[0].Name)
	CheckUnauthorizedStatus(t, resp)

}

func TestUpdatePreferences(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	th.LoginBasic()
	user1 := th.BasicUser

	category := model.NewId()
	preferences1 := model.Preferences{
		{
			UserId:   user1.ClientId,
			Category: category,
			Name:     model.NewId(),
		},
		{
			UserId:   user1.ClientId,
			Category: category,
			Name:     model.NewId(),
		},
		{
			UserId:   user1.ClientId,
			Category: model.NewId(),
			Name:     model.NewId(),
		},
	}

	_, resp := Client.UpdatePreferences(user1.ClientId, &preferences1)
	CheckNoError(t, resp)

	preferences := model.Preferences{
		{
			UserId:   model.NewId(),
			Category: category,
			Name:     model.NewId(),
		},
	}

	_, resp = Client.UpdatePreferences(user1.ClientId, &preferences)
	CheckForbiddenStatus(t, resp)

	preferences = model.Preferences{
		{
			UserId: user1.ClientId,
			Name:   model.NewId(),
		},
	}

	_, resp = Client.UpdatePreferences(user1.ClientId, &preferences)
	CheckBadRequestStatus(t, resp)

	_, resp = Client.UpdatePreferences(th.BasicUser2.ClientId, &preferences)
	CheckForbiddenStatus(t, resp)

	Client.Logout()
	_, resp = Client.UpdatePreferences(user1.ClientId, &preferences1)
	CheckUnauthorizedStatus(t, resp)
}

func TestUpdatePreferencesWebsocket(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	WebSocketClient, err := th.CreateWebSocketClient()
	if err != nil {
		t.Fatal(err)
	}

	WebSocketClient.Listen()
	time.Sleep(300 * time.Millisecond)
	if resp := <-WebSocketClient.ResponseChannel; resp.Status != model.STATUS_OK {
		t.Fatal("should have responded OK to authentication challenge")
	}

	userId := th.BasicUser.ClientId
	preferences := &model.Preferences{
		{
			UserId:   userId,
			Category: model.NewId(),
			Name:     model.NewId(),
		},
		{
			UserId:   userId,
			Category: model.NewId(),
			Name:     model.NewId(),
		},
	}
	_, resp := th.Client.UpdatePreferences(userId, preferences)
	CheckNoError(t, resp)

	timeout := time.After(300 * time.Millisecond)

	waiting := true
	for waiting {
		select {
		case event := <-WebSocketClient.EventChannel:
			if event.Event != model.WEBSOCKET_EVENT_PREFERENCES_CHANGED {
				// Ignore any other events
				continue
			}

			received, err := model.PreferencesFromJson(strings.NewReader(event.Data["preferences"].(string)))
			if err != nil {
				t.Fatal(err)
			}

			for i, preference := range *preferences {
				if preference.UserId != received[i].UserId || preference.Category != received[i].Category || preference.Name != received[i].Name {
					t.Fatal("received incorrect preference")
				}
			}

			waiting = false
		case <-timeout:
			t.Fatal("timed out waiting for preference update event")
		}
	}
}

func TestDeletePreferences(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	th.LoginBasic()

	prefs, _ := Client.GetPreferences(th.BasicUser.ClientId)
	originalCount := len(prefs)

	// save 10 preferences
	var preferences model.Preferences
	for i := 0; i < 10; i++ {
		preference := model.Preference{
			UserId:   th.BasicUser.ClientId,
			Category: model.PREFERENCE_CATEGORY_DIRECT_CHANNEL_SHOW,
			Name:     model.NewId(),
		}
		preferences = append(preferences, preference)
	}

	Client.UpdatePreferences(th.BasicUser.ClientId, &preferences)

	// delete 10 preferences
	th.LoginBasic2()

	_, resp := Client.DeletePreferences(th.BasicUser2.ClientId, &preferences)
	CheckForbiddenStatus(t, resp)

	th.LoginBasic()

	_, resp = Client.DeletePreferences(th.BasicUser.ClientId, &preferences)
	CheckNoError(t, resp)

	_, resp = Client.DeletePreferences(th.BasicUser2.ClientId, &preferences)
	CheckForbiddenStatus(t, resp)

	prefs, _ = Client.GetPreferences(th.BasicUser.ClientId)
	if len(prefs) != originalCount {
		t.Fatal("should've deleted preferences")
	}

	Client.Logout()
	_, resp = Client.DeletePreferences(th.BasicUser.ClientId, &preferences)
	CheckUnauthorizedStatus(t, resp)
}

func TestDeletePreferencesWebsocket(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	userId := th.BasicUser.ClientId
	preferences := &model.Preferences{
		{
			UserId:   userId,
			Category: model.NewId(),
			Name:     model.NewId(),
		},
		{
			UserId:   userId,
			Category: model.NewId(),
			Name:     model.NewId(),
		},
	}
	_, resp := th.Client.UpdatePreferences(userId, preferences)
	CheckNoError(t, resp)

	WebSocketClient, err := th.CreateWebSocketClient()
	if err != nil {
		t.Fatal(err)
	}

	WebSocketClient.Listen()
	time.Sleep(300 * time.Millisecond)
	if resp := <-WebSocketClient.ResponseChannel; resp.Status != model.STATUS_OK {
		t.Fatal("should have responded OK to authentication challenge")
	}

	_, resp = th.Client.DeletePreferences(userId, preferences)
	CheckNoError(t, resp)

	timeout := time.After(30000 * time.Millisecond)

	waiting := true
	for waiting {
		select {
		case event := <-WebSocketClient.EventChannel:
			if event.Event != model.WEBSOCKET_EVENT_PREFERENCES_DELETED {
				// Ignore any other events
				continue
			}

			received, err := model.PreferencesFromJson(strings.NewReader(event.Data["preferences"].(string)))
			if err != nil {
				t.Fatal(err)
			}

			for i, preference := range *preferences {
				if preference.UserId != received[i].UserId || preference.Category != received[i].Category || preference.Name != received[i].Name {
					t.Fatal("received incorrect preference")
				}
			}

			waiting = false
		case <-timeout:
			t.Fatal("timed out waiting for preference delete event")
		}
	}
}
