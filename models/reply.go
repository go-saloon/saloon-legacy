// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Reply struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	AuthorID  uuid.UUID `json:"author_id" db:"author_id"`
	TopicID   uuid.UUID `json:"topic_id" db:"topic_id"`
	Content   string    `json:"content" db:"content"`
	Deleted   bool      `json:"deleted" db:"deleted"`

	Author *User  `json:"-" db:"-"`
	Topic  *Topic `json:"-" db:"-"`
}

type Replies []Reply

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *Reply) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Content, Name: "Content"},
	), nil
}
