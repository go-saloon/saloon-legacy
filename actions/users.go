// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
)

// UserRegisterGet displays a register form
func UsersRegisterGet(c buffalo.Context) error {
	// Make user available inside the html template
	c.Set("user", &models.User{})
	return c.Render(200, r.HTML("users/register.html"))
}

// UsersRegisterPost adds a User to the DB. This function is mapped to the
// path POST /accounts/register
func UsersRegisterPost(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}
	// Bind user to the html form elements
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := user.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make user available inside the html template
		c.Set("user", user)
		// Make the errors available inside the html template
		c.Set("errors", verrs.Errors)
		// Render again the register.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("users/register.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Account created successfully.")
	// and redirect to the home page
	return c.Redirect(302, "/")
}

// UsersLoginGet displays a login form
func UsersLoginGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/login"))
}

// UsersLoginPost logs in a user.
func UsersLoginPost(c buffalo.Context) error {
	user := &models.User{}
	// Bind the user to the html form elements
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}
	tx := c.Value("tx").(*pop.Connection)
	err := user.Authorize(tx)
	if err != nil {
		c.Set("user", user)
		verrs := validate.NewErrors()
		verrs.Add("Login", "Invalid email or password.")
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("users/login"))
	}
	c.Session().Set("current_user_id", user.ID)
	c.Flash().Add("success", "Welcome back!")
	return c.Redirect(302, "/")
}

// UsersLogout clears the session and logs out the user.
func UsersLogout(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "Goodbye!")
	return c.Redirect(302, "/")
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// USerSettingsGet displays the user's informations
func UsersSettingsGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/settings"))
}

func UsersShow(c buffalo.Context) error {
	user := &models.User{}
	uid := c.Param("uid")
	if uid == "" {
		uid = c.Session().Get("current_user_id").(string)
	}
	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Find(user, uid); err != nil {
		return errors.WithStack(err)
	}

	c.Set("user", user)
	return c.Render(200, r.HTML("users/show"))
}

// AdminRequired requires a user to be logged in and to be an admin before accessing a route.
func AdminRequired(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user, ok := c.Value("current_user").(*models.User)
		if ok && user.Admin {
			return next(c)
		}
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/")
	}
}

// UserRequired requires a user to be logged in and to be an admin before accessing a route.
func UserRequired(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user, ok := c.Value("current_user").(*models.User)
		if ok && user != nil {
			return next(c)
		}
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(302, "/")
	}
}
