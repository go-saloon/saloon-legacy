// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grifts

import (
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/pop"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("reset", "Reset database")
	grift.Add("reset", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			return tx.TruncateAll()
		})
	})
})
