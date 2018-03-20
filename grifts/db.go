// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grifts

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"

	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/go-saloon/saloon/actions"
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/pop"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("reset", "Reset (truncate) the whole database")
	grift.Add("reset", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			return tx.TruncateAll()
		})
	})

	grift.Desc("setup", "Initial db setup")
	grift.Add("setup", func(c *grift.Context) error {
		for _, name := range []string{
			"db:setup:create-admin",
			"db:setup:create-forum",
		} {
			err := grift.Run(name, c)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	grift.Desc("setup:create-admin", "Create a default admin account")
	grift.Add("setup:create-admin", func(c *grift.Context) error {
		fset := flag.NewFlagSet("create-admin", flag.ExitOnError)
		pass := fset.String("p", "", "admin password")
		mail := fset.String("e", "", "admin email")

		if *pass == "" {
			*pass = "admin"
		}

		return models.DB.Transaction(func(tx *pop.Connection) error {
			usr := &models.User{
				Username: "admin",
				Email:    *mail,
				Password: *pass,
				Admin:    true,
			}
			pwd, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("could not generate password hash: %v", err)
			}
			usr.PasswordHash = string(pwd)
			usr.Avatar, err = actions.GenAvatar(usr.Username)
			if err != nil {
				return errors.WithStack(err)
			}
			return tx.Create(usr)
		})
	})

	grift.Desc("setup:create-forum", "Create forum welcome message")
	grift.Add("setup:create-forum", func(c *grift.Context) error {
		fset := flag.NewFlagSet("create-forum", flag.ExitOnError)
		fname := fset.String("logo", "", "path to a logo image")
		title := fset.String("title", "Saloon", "title of the forum")
		descr := fset.String("descr", "A nice forum", "description of the forum")

		err := fset.Parse(c.Args)
		if err != nil {
			return errors.WithStack(err)
		}

		return models.DB.Transaction(func(tx *pop.Connection) error {
			var (
				err  error
				logo []byte
			)

			if *fname != "" {
				f, err := os.Open(*fname)
				if err != nil {
					return errors.WithStack(err)
				}
				defer f.Close()
				img, _, err := image.Decode(f)
				if err != nil {
					return errors.WithStack(err)
				}
				buf := new(bytes.Buffer)
				err = png.Encode(buf, img)
				if err != nil {
					return errors.WithStack(err)
				}
				logo = buf.Bytes()
			}
			forum := &models.Forum{
				Title:       *title,
				Description: *descr,
				Logo:        logo,
			}
			err = tx.Create(forum)
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
	})

	grift.Desc("update-user-password", "Update and change a user password")
	grift.Add("update-user-password", func(c *grift.Context) error {
		fset := flag.NewFlagSet("update-user-password", flag.ExitOnError)
		usrname := fset.String("u", "", "user name")
		password := fset.String("p", "", "user new password")

		if err := fset.Parse(c.Args); err != nil {
			return errors.WithStack(err)
		}

		return models.DB.Transaction(func(tx *pop.Connection) error {
			usr := new(models.User)
			usr.Username = *usrname
			if err := tx.Where("username = ?", *usrname).First(usr); err != nil {
				return errors.WithStack(err)
			}

			usr.Password = *password
			pwd, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
			if err != nil {
				return errors.WithStack(err)
			}
			usr.PasswordHash = string(pwd)

			if err := tx.Update(usr); err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
	})
})
