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

/*
List-Post: <mailto:reply+00105748c619555d4a6c80b4faccec22003b863b33e73ae092cf0000000116c2ac9c92a169ce1238ebbe@reply.github.com>
List-Unsubscribe: <mailto:unsub+00105748c619555d4a6c80b4faccec22003b863b33e73ae092cf0000000116c2ac9c92a169ce1238ebbe@reply.github.com>, <https://github.com/notifications/unsubscribe/ABBXSLhgVLtfNtdMGG1Y0aRw9bFiNJc_ks5teuIcgaJpZM4Ss4xE>
*/

func NewTopicNotify(c buffalo.Context, topic *models.Topic, recpts []models.User) error {
	m := mail.NewMessage()
	m.SetHeader("Reply-To", notify.ReplyTo)
	m.SetHeader("Message-ID", fmt.Sprintf("<topic/%s@%s>", topic.ID, notify.MessageID))
	m.SetHeader("List-ID", notify.ListID)
	m.SetHeader("List-Archive", notify.ListArchive)
	m.SetHeader("List-Unsubscribe", notify.ListUnsubscribe)
	m.SetHeader("X-Auto-Response-Suppress", "All")

	m.Subject = notify.SubjectHdr + " " + topic.Title
	m.From = fmt.Sprintf("%s <%s>", topic.Author.Username, notify.From)
	m.To = nil
	m.Bcc = nil
	for _, usr := range recpts {
		m.Bcc = append(m.Bcc, usr.Email)
	}

	data := map[string]interface{}{
		"content":     topic.Content,
		"unsubscribe": notify.ListUnsubscribe,
	}

	err := m.AddBodies(
		data,
		r.Plain("mail/notify.txt"),
		r.HTML("mail/notify.html"),
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
	m := mail.NewMessage()
	m.SetHeader("Reply-To", notify.ReplyTo)
	m.SetHeader("Message-ID", fmt.Sprintf("<topic/%s/%s@%s>", topic.ID, reply.ID, notify.MessageID))
	m.SetHeader("In-Reply-To", fmt.Sprintf("<topic/%s@%s>", topic.ID, notify.InReplyTo))
	m.SetHeader("List-ID", notify.ListID)
	m.SetHeader("List-Archive", notify.ListArchive)
	m.SetHeader("List-Unsubscribe", notify.ListUnsubscribe)
	m.SetHeader("X-Auto-Response-Suppress", "All")

	m.Subject = notify.SubjectHdr + " " + topic.Title
	m.From = fmt.Sprintf("%s <%s>", reply.Author.Username, notify.From)
	m.To = nil
	m.Bcc = nil
	for _, usr := range recpts {
		m.Bcc = append(m.Bcc, usr.Email)
	}

	data := map[string]interface{}{
		"content":     reply.Content,
		"unsubscribe": notify.ListUnsubscribe,
	}

	err := m.AddBodies(
		data,
		r.Plain("mail/notify.txt"),
		r.HTML("mail/notify.html"),
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
