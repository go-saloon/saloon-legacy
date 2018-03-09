package actions

import "github.com/gobuffalo/buffalo"

// RepliesCreate default implementation.
func RepliesCreate(c buffalo.Context) error {
	return c.Render(200, r.HTML("replies/create.html"))
}

// RepliesEdit default implementation.
func RepliesEdit(c buffalo.Context) error {
	return c.Render(200, r.HTML("replies/edit.html"))
}

// RepliesDelete default implementation.
func RepliesDelete(c buffalo.Context) error {
	return c.Render(200, r.HTML("replies/delete.html"))
}
