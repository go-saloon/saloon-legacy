package actions

import "github.com/gobuffalo/buffalo"

// TopicsIndex default implementation.
func TopicsIndex(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/index.html"))
}

func TopicsCreateGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/create.html"))
}

func TopicsCreatePost(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/create.html"))
}

// TopicsEdit default implementation.
func TopicsEdit(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/edit.html"))
}

// TopicsDelete default implementation.
func TopicsDelete(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/delete.html"))
}

// TopicsDetail default implementation.
func TopicsDetail(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/detail.html"))
}
