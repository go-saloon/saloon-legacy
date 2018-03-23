// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"sort"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/disintegration/letteravatar"
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
	avatar, err := GenAvatar(user.Username)
	if err != nil {
		return errors.WithStack(err)
	}
	user.Avatar = avatar
	tx := c.Value("tx").(*pop.Connection)

	// subscribe user with all categories.
	// FIXME(sbinet) we should make the list of default categories
	// customizable at the application level...
	// see:
	//   go-saloon/saloon#8
	cats := new(models.Categories)
	if err := tx.All(cats); err != nil {
		return errors.WithStack(err)
	}
	for _, cat := range *cats {
		user.AddSubscription(cat.ID)
	}

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
	c.Session().Set("current_user_id", user.ID)
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
		verrs.Add("Login", "Invalid user or password.")
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

// UserSettings displays the user's informations
func UsersSettings(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	cats := new(models.Categories)
	if err := tx.All(cats); err != nil {
		return errors.WithStack(err)
	}
	sort.Sort(cats)
	c.Set("categories", cats)
	c.Set("avatar", new(models.Avatar))
	usr := c.Value("current_user").(*models.User)
	if usr.Admin {
		users := new(models.Users)
		if err := tx.All(users); err != nil {
			return errors.WithStack(err)
		}
		sort.Sort(users)
		c.Set("users", users)
	}
	return c.Render(200, r.HTML("users/settings"))
}

func UsersSettingsAddSubscription(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	cat := new(models.Category)
	if err := tx.Find(cat, c.Param("cid")); err != nil {
		return errors.WithStack(err)
	}
	usr := new(models.User)
	if err := tx.Find(usr, c.Param("uid")); err != nil {
		return errors.WithStack(err)
	}
	usr.AddSubscription(cat.ID)
	cat.AddSubscriber(usr.ID)

	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}
	if err := tx.Update(cat); err != nil {
		return errors.WithStack(err)
	}
	return c.Redirect(302, "/users/settings")
}

func UsersSettingsRemoveSubscription(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	cat := new(models.Category)
	if err := tx.Find(cat, c.Param("cid")); err != nil {
		return errors.WithStack(err)
	}
	usr := new(models.User)
	if err := tx.Find(usr, c.Param("uid")); err != nil {
		return errors.WithStack(err)
	}
	usr.RemoveSubscription(cat.ID)
	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}
	cat.RemoveSubscriber(usr.ID)
	if err := tx.Update(cat); err != nil {
		return errors.WithStack(err)
	}
	return c.Redirect(302, "/users/settings")
}

func UsersSettingsUpdateAvatar(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	cats := new(models.Categories)
	if err := tx.All(cats); err != nil {
		return errors.WithStack(err)
	}
	sort.Sort(cats)

	usr := c.Value("current_user").(*models.User)
	f, err := c.File("avatar")
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	if !f.Valid() {
		verrs := validate.NewErrors()
		verrs.Add("Avatar Upload", "Invalid file")
		c.Set("errors", verrs.Errors)
		c.Set("avatar", new(models.Avatar))
		return c.Render(422, r.HTML("users/settings"))
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return errors.WithStack(err)
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return errors.WithStack(err)
	}

	usr.Avatar = buf.Bytes()
	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/users/settings")
}

func UsersSettingsUpdateName(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	usr := c.Value("current_user").(*models.User)
	if err := c.Bind(usr); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/users/settings")
}

func UsersSettingsUpdateEmail(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	usr := c.Value("current_user").(*models.User)
	if err := c.Bind(usr); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}
	return c.Redirect(302, "/users/settings")
}

func UsersSettingsUpdatePassword(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	usr := c.Value("current_user").(*models.User)
	if err := c.Bind(usr); err != nil {
		return errors.WithStack(err)
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	usr.PasswordHash = string(pwd)

	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}
	return c.Redirect(302, "/users/settings")
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

func GenAvatar(name string) ([]byte, error) {
	const avatarSize = 100
	letter, _ := utf8.DecodeRuneInString(name)
	img, err := letteravatar.Draw(avatarSize, unicode.ToUpper(letter), nil)
	if err != nil {
		return nil, fmt.Errorf("could not generate letteravatar: %v", err)
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, fmt.Errorf("could not encode letteravatar to PNG: %v", err)
	}
	return buf.Bytes(), nil
}
