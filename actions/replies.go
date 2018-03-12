package actions

import (
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func RepliesCreateGet(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	reply := &models.Reply{}
	topic := &models.Topic{}
	if err := tx.Find(topic, c.Param("tid")); err != nil {
		return c.Error(404, err)
	}
	c.Set("reply", reply)
	c.Set("topic", topic)
	reply.TopicID = topic.ID
	return c.Render(200, r.HTML("replies/create.html"))
}

func RepliesCreatePost(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	reply := new(models.Reply)
	topic := new(models.Topic)
	user := c.Value("current_user").(*models.User)
	if err := c.Bind(reply); err != nil {
		return errors.WithStack(err)
	}
	if err := tx.Find(topic, c.Param("tid")); err != nil {
		return c.Error(404, err)
	}
	c.Set("topic", topic)
	reply.AuthorID = user.ID
	reply.TopicID = topic.ID

	verrs, err := tx.ValidateAndCreate(reply)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("reply", reply)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("replies/create"))
	}
	c.Flash().Add("success", "New reply added successfully.")
	return c.Redirect(302, "/topics/detail/%s", topic.ID)
}

// RepliesEdit default implementation.
func RepliesEdit(c buffalo.Context) error {
	return c.Render(200, r.HTML("replies/edit.html"))
}

// RepliesDelete default implementation.
func RepliesDelete(c buffalo.Context) error {
	return c.Render(200, r.HTML("replies/delete.html"))
}

func RepliesDetail(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	topic := new(models.Topic)
	if err := tx.Find(topic, c.Param("tid")); err != nil {
		return c.Error(404, err)
	}
	replies := new(models.Replies)
	if err := tx.BelongsTo(topic).All(replies); err != nil {
		return c.Error(404, err)
	}
	c.Set("topic", topic)
	c.Set("replies", replies)
	return c.Render(200, r.HTML("replies/detail"))
}
