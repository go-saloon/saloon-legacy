// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models_test

import "github.com/go-saloon/saloon/models"

func (ms *ModelSuite) Test_User_Create() {
	n, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, n)

	u := &models.User{
		Username:        "usr1",
		Email:           "user@example.com",
		Password:        "password",
		PasswordConfirm: "password",
	}
	ms.Zero(u.PasswordHash)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.NotZero(u.PasswordHash)
	ms.False(u.Admin)

	n, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(1, n)
}

func (ms *ModelSuite) Test_User_Create_ValidationErrors() {
	n, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, n)

	u := &models.User{
		Email:           "user@example.com",
		Password:        "password",
		PasswordConfirm: "password",
	}
	ms.Zero(u.PasswordHash)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())

	n, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(0, n)
}
