// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grifts

import (
	. "github.com/markbates/grift/grift"
)

var _ = Namespace("db", func() {

	Namespace("setup", func() {

		Desc("create_admin", "Create a default admin account")
		Add("create_admin", func(c *Context) error {
			return nil
		})

	})

})
