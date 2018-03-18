// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mailers

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"
	"net"
	"text/template"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

var smtp mail.Sender
var r *render.Engine

var notify struct {
	ReplyTo    string
	MessageID  string
	InReplyTo  string
	ListID     string
	SubjectHdr string
	From       string
}

func init() {

	// Pulling config from the env.
	send := envy.Get("SALOON_SEND_MAIL", "")
	port := envy.Get("SMTP_PORT", "1025")
	host := envy.Get("SMTP_HOST", "localhost")
	user := envy.Get("SMTP_USER", "")
	password := envy.Get("SMTP_PASSWORD", "")

	notify.ReplyTo = envy.Get("SALOON_MAIL_NOTIFY_REPLY_TO", "")
	notify.MessageID = envy.Get("SALOON_MAIL_NOTIFY_MESSAGE_ID", "")
	notify.InReplyTo = envy.Get("SALOON_MAIL_NOTIFY_IN_REPLY_TO", "")
	notify.ListID = envy.Get("SALOON_MAIL_NOTIFY_LIST_ID", "")
	notify.SubjectHdr = envy.Get("SALOON_MAIL_NOTIFY_SUBJECT_HDR", "")
	notify.From = envy.Get("SALOON_MAIL_NOTIFY_FROM", "")

	var err error

	switch send {
	case "1", "y", "yes", "Y":
		smtp, err = mail.NewSMTPSender(host, port, user, password)
	case "remote":
		smtp, err = newRelaySender(host, port, user, password)
	default:
		smtp = noMailSender{}
	}
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

// noMailSender is a mail.Sender with no actual delivery.
type noMailSender struct{}

func (noMailSender) Send(mail.Message) error {
	return nil
}

type relaySender struct {
	addr string
}

func newRelaySender(host, port, user, password string) (*relaySender, error) {
	return &relaySender{host + ":" + port}, nil
}

func (rs relaySender) Send(m mail.Message) error {
	conn, err := net.Dial("tcp", rs.addr)
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	o := new(bytes.Buffer)
	for i := range m.Attachments {
		att := &m.Attachments[i]
		r := new(bytes.Buffer)
		_, err := io.Copy(r, att.Reader)
		if err != nil {
			return errors.WithStack(err)
		}
		att.Reader = r
	}

	err = gob.NewEncoder(o).Encode(m)
	if err != nil {
		return errors.WithStack(err)
	}

	var hdr [8]byte
	n := int64(o.Len())
	binary.BigEndian.PutUint64(hdr[:], uint64(n))
	_, err = conn.Write(hdr[:])
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = conn.Write(o.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
