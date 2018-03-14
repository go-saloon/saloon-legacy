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
	reply.Topic = topic
	reply.Author = c.Value("current_user").(*models.User)
	return c.Render(200, r.HTML("replies/create.html"))
}

func RepliesCreatePost(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	reply := new(models.Reply)
	user := c.Value("current_user").(*models.User)
	if err := c.Bind(reply); err != nil {
		return errors.WithStack(err)
	}
	topic, err := loadTopic(c, c.Param("tid"))
	if err != nil {
		return c.Error(404, err)
	}
	c.Set("topic", topic)
	reply.AuthorID = user.ID
	reply.Author = user
	reply.TopicID = topic.ID
	reply.Topic = topic

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

	err = notifyTopic(c, topic, reply)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, "/topics/detail/%s", topic.ID)
}

// RepliesEdit default implementation.
func RepliesEdit(c buffalo.Context) error {
	return c.Render(200, r.HTML("replies/edit.html"))
}

func RepliesDelete(c buffalo.Context) error {
	reply, err := loadReply(c, c.Param("rid"))
	if err != nil {
		return errors.WithStack(err)
	}
	usr := c.Value("current_user").(*models.User)
	if !usr.Admin && usr.ID != reply.AuthorID {
		c.Flash().Add("danger", "You are not authorized to delete this reply")
		return c.Redirect(302, "/topics/detail/%s", reply.TopicID)
	}
	tx := c.Value("tx").(*pop.Connection)
	reply.Deleted = true
	if err := tx.Update(reply); err != nil {
		return errors.WithStack(err)
	}
	c.Flash().Add("success", "Reply deleted successfuly.")
	return c.Redirect(302, "/topics/detail/%s", reply.TopicID)
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

func loadReply(c buffalo.Context, id string) (*models.Reply, error) {
	tx := c.Value("tx").(*pop.Connection)
	reply := &models.Reply{}
	if err := c.Bind(reply); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := tx.Find(reply, id); err != nil {
		return nil, c.Error(404, err)
	}
	topic := new(models.Topic)
	if err := tx.Find(topic, reply.TopicID); err != nil {
		return nil, c.Error(404, err)
	}
	usr := new(models.User)
	if err := tx.Find(usr, reply.AuthorID); err != nil {
		return nil, c.Error(404, err)
	}
	reply.Topic = topic
	reply.Author = usr
	return reply, nil
}
