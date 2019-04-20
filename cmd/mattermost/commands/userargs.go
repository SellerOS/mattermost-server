// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"github.com/mattermost/mattermost-server/app"
	"github.com/mattermost/mattermost-server/model"
)

func getUsersFromUserArgs(a *app.App, userArgs []string) []*model.UserIms {
	users := make([]*model.UserIms, 0, len(userArgs))
	for _, userArg := range userArgs {
		user := getUserFromUserArg(a, userArg)
		users = append(users, user)
	}
	return users
}

func getUserFromUserArg(a *app.App, userArg string) *model.UserIms {
	var user *model.UserIms
	if result := <-a.Srv.Store.User().GetByEmail(userArg); result.Err == nil {
		user = result.Data.(*model.UserIms)
	}

	if user == nil {
		if result := <-a.Srv.Store.User().GetByUsername(userArg); result.Err == nil {
			user = result.Data.(*model.UserIms)
		}
	}

	if user == nil {
		user, _ = a.Srv.Store.User().Get(userArg)
	}

	return user
}
