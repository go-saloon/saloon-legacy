// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mailers

import (
	"bytes"
	"log"
	"text/template"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
)

var smtp mail.Sender
var r *render.Engine

func init() {

	// Pulling config from the env.
	port := envy.Get("SMTP_PORT", "1025")
	host := envy.Get("SMTP_HOST", "localhost")
	user := envy.Get("SMTP_USER", "")
	password := envy.Get("SMTP_PASSWORD", "")

	var err error

	smtp, err = mail.NewSMTPSender(host, port, user, password)
	if err != nil {
		log.Fatal(err)
	}

	r = render.New(render.Options{
		HTMLLayout:   "mail/layout.html",
		TemplatesBox: packr.NewBox("../templates"),
		Helpers:      render.Helpers{},
		TemplateEngines: map[string]render.TemplateEngine{
			"txt": PlainTextTemplateEngine,
		},
	})

}

func PlainTextTemplateEngine(input string, data map[string]interface{}, helpers map[string]interface{}) (string, error) {
	// since go templates don't have the concept of an optional map argument like Plush does
	// add this "null" map so it can be used in templates like this:
	// {{ partial "flash.html" .nilOpts }}
	data["nilOpts"] = map[string]interface{}{}

	t := template.New(input)
	if helpers != nil {
		t = t.Funcs(helpers)
	}

	t, err := t.Parse(input)
	if err != nil {
		return "", err
	}

	bb := &bytes.Buffer{}
	err = t.Execute(bb, data)
	return bb.String(), err
}
