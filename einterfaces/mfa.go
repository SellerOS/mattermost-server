// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package einterfaces

import (
	"github.com/mattermost/mattermost-server/model"
)

type MfaInterface interface {
	GenerateSecret(user *model.UserInfo) (string, []byte, *model.AppError)
	Activate(user *model.UserIms, token string) *model.AppError
	Deactivate(userId string) *model.AppError
	ValidateToken(secret, token string) (bool, *model.AppError)
}
