// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// SetCurrentForum attempts to find a forum definition in the db.
func SetCurrentForum(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tx := c.Value("tx").(*pop.Connection)
		forum := &models.Forum{}
		err := tx.First(forum)
		if err != nil {
			return errors.WithStack(err)
		}
		c.Set("forum", forum)
		return next(c)
	}
}
