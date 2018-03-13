// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Forum struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at,utc"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at,utc"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	BaseAddr    string    `json:"base_addr" db:"base_addr"`
}

// String is not required by pop and may be deleted
func (f Forum) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Forums is not required by pop and may be deleted
type Forums []Forum

// String is not required by pop and may be deleted
func (f Forums) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (f *Forum) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: f.BaseAddr, Name: "BaseAddr"},
		&validators.StringIsPresent{Field: f.Title, Name: "Title"},
		&validators.StringIsPresent{Field: f.Description, Name: "Description"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (f *Forum) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (f *Forum) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type validURL struct {
	Name  string
	Field string
}

func (v *validURL) IsValid(errors *validate.Errors) {
	_, err := url.Parse(v.Field)
	if err != nil {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("The address %q is not a valid URL", v.Field))
	}
}
