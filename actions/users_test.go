// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

/* FIXME(sbinet)

import "github.com/go-saloon/saloon/models"

func (as *ActionSuite) Test_Users_Register() {
	res := as.HTML("/users/register").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_Users_Create() {
	n, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, n)

	u := &models.User{
		Username:        "usr1",
		Email:           "user@example.com",
		Password:        "password",
		PasswordConfirm: "password",
	}

	res := as.HTML("/users/register").Post(u)
	as.Equal(302, res.Code)

	n, err = as.DB.Count("users")
	as.NoError(err)
	as.Equal(1, n)
}
*/
