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
	topic.Author = c.Value("current_user").(*models.User)
	if err := c.Bind(topic); err != nil {
		return errors.WithStack(err)
	}
	cat := new(models.Category)
	if err := tx.Find(cat, c.Param("cid")); err != nil {
		return c.Error(404, err)
	}
	topic.Category = cat
	topic.AuthorID = topic.Author.ID
	topic.CategoryID = topic.Category.ID
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
	topic, err := loadTopic(c, c.Param("tid"))
	if err != nil {
		return errors.WithStack(err)
	}
	c.Set("topic", topic)
	c.Set("category", topic.Category)
	c.Set("replies", &topic.Replies)
	return c.Render(200, r.HTML("topics/detail"))
}

func loadTopic(c buffalo.Context, tid string) (*models.Topic, error) {
	tx := c.Value("tx").(*pop.Connection)
	topic := &models.Topic{}
	if err := c.Bind(topic); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := tx.Find(topic, tid); err != nil {
		return nil, c.Error(404, err)
	}
	cat := new(models.Category)
	if err := tx.Find(cat, topic.CategoryID); err != nil {
		return nil, c.Error(404, err)
	}
	usr := new(models.User)
	if err := tx.Find(usr, topic.AuthorID); err != nil {
		return nil, c.Error(404, err)
	}
	if err := tx.BelongsTo(topic).All(&topic.Replies); err != nil {
		return nil, c.Error(404, err)
	}
	topic.Category = cat
	topic.Author = usr
	for i := range topic.Replies {
		reply, err := loadReply(c, topic.Replies[i].ID.String())
		if err != nil {
			return nil, c.Error(404, err)
		}
		topic.Replies[i] = *reply
	}
	return topic, nil
}
