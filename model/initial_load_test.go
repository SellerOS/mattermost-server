// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"strings"
	"testing"
)

func TestInitialLoadJson(t *testing.T) {
	u := &User{ClientId: NewId()}
	o := InitialLoad{User: u}
	json := o.ToJson()
	ro := InitialLoadFromJson(strings.NewReader(json))

	if o.User.ClientId != ro.User.ClientId {
		t.Fatal("Ids do not match")
	}
}
