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
	topic.Author = user
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
func TopicsEditGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/edit.html"))
}

// TopicsEdit default implementation.
func TopicsEditPost(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/edit.html"))
}

// TopicsDelete default implementation.
func TopicsDeleteGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("topics/delete.html"))
}

// TopicsDelete default implementation.
func TopicsDeletePost(c buffalo.Context) error {
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
	cat := new(models.Category)
	if err := tx.Find(cat, topic.CategoryID); err != nil {
		return c.Error(404, err)
	}
	c.Set("category", cat)
	author := new(models.User)
	if err := tx.Find(author, topic.AuthorID); err != nil {
		return c.Error(404, err)
	}
	topic.Author = author
	q := tx.PaginateFromParams(c.Params())
	replies := new(models.Replies)
	if err := q.BelongsTo(topic).All(replies); err != nil {
		return c.Error(404, err)
	}
	topic.Replies = *replies
	c.Set("replies", replies)
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("topics/detail"))
}
