// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grifts

import (
	. "github.com/markbates/grift/grift"
)

var _ = Namespace("db", func() {

	Desc("reset", "Reset database")
	Add("reset", func(c *Context) error {
		return nil
	})

})
