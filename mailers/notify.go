// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mailers

import (
	"fmt"

	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/mail"
	"github.com/pkg/errors"
)

func NewTopicNotify(c buffalo.Context, topic *models.Topic, recpts []models.User) error {
	// Creates a new message
	m := mail.NewMessage()
	m.SetHeader("Reply-To", notify.ReplyTo)
	m.SetHeader("Message-ID", fmt.Sprintf("<topic/%s@%s>", topic.ID, notify.MessageID))
	m.SetHeader("List-ID", notify.ListID)

	m.Subject = notify.SubjectHdr + " " + topic.Title
	m.From = fmt.Sprintf("%s <%s>", topic.Author.Username, notify.From)
	m.To = nil
	m.Bcc = nil
	for _, usr := range recpts {
		m.Bcc = append(m.Bcc, usr.Email)
	}

	// Data that will be used inside the templates when rendering.
	data := map[string]interface{}{
		"content": topic.Content,
	}

	// You can add multiple bodies to the message you're creating to have content-types alternatives.
	err := m.AddBodies(
		data,
		r.HTML("mail/notify.html"),
		r.Plain("mail/notify.txt"),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	err = smtp.Send(m)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func NewReplyNotify(c buffalo.Context, topic *models.Topic, reply *models.Reply, recpts []models.User) error {
	// Creates a new message
	m := mail.NewMessage()
	m.SetHeader("Reply-To", notify.ReplyTo)
	m.SetHeader("Message-ID", fmt.Sprintf("<topic/%s/%s@%s>", topic.ID, reply.ID, notify.MessageID))
	m.SetHeader("In-Reply-To", fmt.Sprintf("<topic/%s@%s>", topic.ID, notify.InReplyTo))
	m.SetHeader("List-ID", notify.ListID)

	m.Subject = notify.SubjectHdr + " " + topic.Title
	m.From = fmt.Sprintf("%s <%s>", reply.Author.Username, notify.From)
	m.To = nil
	m.Bcc = nil
	for _, usr := range recpts {
		m.Bcc = append(m.Bcc, usr.Email)
	}

	// Data that will be used inside the templates when rendering.
	data := map[string]interface{}{
		"content": reply.Content,
	}

	// You can add multiple bodies to the message you're creating to have content-types alternatives.
	err := m.AddBodies(
		data,
		r.HTML("mail/notify.html"),
		r.Plain("mail/notify.txt"),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	err = smtp.Send(m)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
