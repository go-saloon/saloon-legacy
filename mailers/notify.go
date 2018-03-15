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
Transport; Wed, 14 Mar 2018 17:30:38 +0000
Received: from localhost.localdomain (188.184.69.125) by cernmx.cern.ch (188.184.36.24) with Microsoft SMTP Server id 14.3.319.2; Wed, 14 Mar 2018 18:30:21 +0100
Return-Path: replies+verp-435445b6111626719fedf95e46d81dc9@root-forum-mail.cern.ch
Date: Wed, 14 Mar 2018 17:30:20 +0000
From: Viesturs <root.discourse@cern.ch>
Reply-To: ROOT Forum <replies+ee5527776516bff82feb239da4e449ca@root-forum-mail.cern.ch>
To: <seb.binet@gmail.com>
Message-ID: <topic/28306/124700@root-forum.cern.ch>
In-Reply-To: <topic/28306/124690@root-forum.cern.ch>
References: <topic/28306@root-forum.cern.ch> <topic/28306/124690@root-forum.cern.ch>
Subject: [ROOT Forum] [ROOT] Load a standalone sharedlibrary
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="--==_mimepart_5aa95c2c8115b_1ff3fd51e9899e05662fb"; charset="UTF-8"
Content-Transfer-Encoding: 7bit
List-Unsubscribe: <http://root-forum.cern.ch/email/unsubscribe/955b9d9ea7d609b54171b3b7e0ca582dd6dc3832c88860999c80c9e5eae5fb62>
X-Auto-Response-Suppress: All
Auto-Submitted: auto-generated
Precedence: list
List-ID: <root.root-forum.cern.ch>
List-Archive: http://root-forum.cern.ch/t/load-a-standalone-sharedlibrary/28306
*/

func NewTopicNotify(c buffalo.Context, topic *models.Topic, recpts []models.User) error {
	// Creates a new message
	m := mail.NewMessage()
	m.SetHeader("Message-ID", fmt.Sprintf("<topic/%s@saloon.com>", topic.ID))
	m.SetHeader("List-ID", "forum.saloon.com")

	m.From = fmt.Sprintf("%s <alternatiba63.forum@gmail.com>", topic.Author.Username)
	m.Subject = topic.Title
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
	m.SetHeader("Message-ID", fmt.Sprintf("<topic/%s/%s@saloon.com>", topic.ID, reply.ID))
	m.SetHeader("In-Reply-To", fmt.Sprintf("<topic/%s@saloon.com>", topic.ID))
	m.SetHeader("List-ID", "forum.saloon.com")

	m.From = fmt.Sprintf("%s <alternatiba63.forum@gmail.com>", reply.Author.Username)
	m.Subject = topic.Title
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
