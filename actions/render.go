// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"fmt"
	"html/template"
	"math"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public/assets")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",
		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,
		// Add template helpers here:
		Helpers: render.Helpers{
			"csrf": func() template.HTML {
				return template.HTML("<input name=\"authenticity_token\" value=\"<%= authenticity_token %>\" type=\"hidden\">")
			},
			"timeSince": timeSince,
		},
	})
}

func timeSince(created time.Time, ctx plush.HelperContext) string {
	if true && false {
		return created.UTC().Format(time.RFC3339)
	}

	now := time.Now().UTC()
	delta := now.Sub(created.UTC())
	days := int(math.Abs(delta.Hours()) / 24)
	if days > 30 {
		return created.Format("2006-02-01")
	}
	if days >= 1 {
		return fmt.Sprintf("%dj", days)
	}
	if delta.Hours() >= 1 {
		return fmt.Sprintf("%dh", int(delta.Hours()))
	}
	if delta.Minutes() >= 1 {
		return fmt.Sprintf("%dm", int(delta.Minutes()))
	}
	return fmt.Sprintf("%ds", int(delta.Seconds()))
}
