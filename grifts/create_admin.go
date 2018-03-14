// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grifts

import (
	"bytes"
	"fmt"
	"image/png"
	"unicode"
	"unicode/utf8"

	"github.com/disintegration/letteravatar"
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("setup", "Initial db setup")
	grift.Add("setup", func(c *grift.Context) error {
		for _, name := range []string{
			"db:setup:create-admin",
			"db:setup:create-dummy-users",
		} {
			err := grift.Run(name, c)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	grift.Namespace("setup", func() {

		grift.Desc("create-admin", "Create a default admin account")
		grift.Add("create-admin", func(c *grift.Context) error {
			return models.DB.Transaction(func(tx *pop.Connection) error {
				usr := &models.User{
					Username: "admin",
					Email:    "admin@example.com",
					Password: "admin",
					Admin:    true,
				}
				pwd, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
				if err != nil {
					return fmt.Errorf("could not generate password hash: %v", err)
				}
				usr.PasswordHash = string(pwd)
				usr.Avatar, err = genAvatar(usr.Username)
				if err != nil {
					return errors.WithStack(err)
				}
				return tx.Create(usr)
			})
		})

		grift.Desc("create-dummy-users", "Create dummy user accounts")
		grift.Add("create-dummy-users", func(c *grift.Context) error {
			return models.DB.Transaction(func(tx *pop.Connection) error {
				for _, usr := range []*models.User{
					{
						Username:  "toto",
						Email:     "toto@example.com",
						FirstName: nulls.String{String: "Mr", Valid: true},
						LastName:  nulls.String{String: "Toto", Valid: true},
					},
					{
						Username:  "tata",
						Email:     "tata@example.com",
						FirstName: nulls.String{String: "Mme", Valid: true},
						LastName:  nulls.String{String: "Tata", Valid: true},
					},
				} {
					usr.Password = usr.Username
					usr.PasswordConfirm = usr.Username
					pwd, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
					if err != nil {
						return fmt.Errorf("could not generate password hash: %v", err)
					}
					usr.PasswordHash = string(pwd)
					usr.Avatar, err = genAvatar(usr.Username)
					if err != nil {
						return errors.WithStack(err)
					}
					err = tx.Create(usr)
					if err != nil {
						return errors.WithStack(err)
					}
				}
				return nil
			})
		})
	})

})

func genAvatar(name string) ([]byte, error) {
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
