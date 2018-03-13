// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"fmt"
	"html/template"
	"math"
	"time"

	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
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
			"userName":      userName,
			"categoryTitle": categoryTitle,
			"topicTitle":    topicTitle,
			"timeSince":     timeSince,
		},
	})
}

func userName(id uuid.UUID, ctx plush.HelperContext) string {
	tx := ctx.Value("tx").(*pop.Connection)
	v := new(models.User)
	if err := tx.Find(v, id); err != nil {
		return "N/A"
	}
	return v.Username
}

func categoryTitle(id uuid.UUID, ctx plush.HelperContext) string {
	tx := ctx.Value("tx").(*pop.Connection)
	v := new(models.Category)
	if err := tx.Find(v, id); err != nil {
		return "N/A"
	}
	return v.Title
}

func topicTitle(id uuid.UUID, ctx plush.HelperContext) string {
	tx := ctx.Value("tx").(*pop.Connection)
	v := new(models.Topic)
	if err := tx.Find(v, id); err != nil {
		return "N/A"
	}
	return v.Title
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
