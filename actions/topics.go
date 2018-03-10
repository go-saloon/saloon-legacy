package actions

import (
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// TopicsIndex default implementation.
func TopicsIndex(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/index.html"))
}

func TopicsCreateGet(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	topic := &models.Topic{}
	c.Set("topic", topic)
	cat := &models.Category{}
	if err := tx.Find(cat, c.Param("cid")); err != nil {
		return c.Error(404, err)
	}
	c.Set("category", cat)
	topic.CategoryID = cat.ID

	return c.Render(200, r.HTML("topics/create"))
}

func TopicsCreatePost(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	topic := &models.Topic{}
	user := c.Value("current_user").(*models.User)
	if err := c.Bind(topic); err != nil {
		return errors.WithStack(err)
	}
	cat := &models.Category{}
	if err := tx.Find(cat, c.Param("cid")); err != nil {
		return c.Error(404, err)
	}
	topic.AuthorID = user.ID
	topic.CategoryID = cat.ID
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(topic)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("topic", topic)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("topics/create"))
	}
	c.Flash().Add("success", "New topic added successfully.")
	return c.Redirect(302, "/topics/detail/%s", topic.ID)
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
	tx := c.Value("tx").(*pop.Connection)
	topic := &models.Topic{}
	if err := c.Bind(topic); err != nil {
		return errors.WithStack(err)
	}
	if err := tx.Find(topic, c.Param("tid")); err != nil {
		return c.Error(404, err)
	}
	c.Set("topic", topic)
	return c.Render(200, r.HTML("topics/detail"))
}
