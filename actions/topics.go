// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"sort"

	"github.com/go-saloon/saloon/mailers"
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
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
	topic.AddSubscriber(topic.AuthorID)
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

	err = newTopicNotify(c, topic)
	if err != nil {
		return errors.WithStack(err)
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
func TopicsDelete(c buffalo.Context) error {
	topic, err := loadTopic(c, c.Param("tid"))
	if err != nil {
		return errors.WithStack(err)
	}
	usr := c.Value("current_user").(*models.User)
	if !usr.Admin && usr.ID != topic.AuthorID {
		c.Flash().Add("danger", "You are not authorized to delete this topic")
		return c.Redirect(302, "/topics/detail/%s", topic.ID)
	}
	tx := c.Value("tx").(*pop.Connection)
	topic.Deleted = true
	if err := tx.Update(topic); err != nil {
		return errors.WithStack(err)
	}
	c.Flash().Add("success", "Topic deleted successfuly.")
	return c.Redirect(302, "/categories/detail/%s", topic.CategoryID)
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

func TopicsAddSubscriber(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	topic, err := loadTopic(c, c.Param("tid"))
	if err != nil {
		return errors.WithStack(err)
	}
	usr := c.Value("current_user").(*models.User)
	topic.AddSubscriber(usr.ID)

	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Update(topic); err != nil {
		return errors.WithStack(err)
	}

	c.Set("topic", topic)
	c.Set("category", topic.Category)
	c.Set("replies", &topic.Replies)
	return c.Redirect(302, "/topics/detail/%s", topic.ID)
}

func TopicsRemoveSubscriber(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	topic, err := loadTopic(c, c.Param("tid"))
	if err != nil {
		return errors.WithStack(err)
	}
	usr := c.Value("current_user").(*models.User)
	topic.RemoveSubscriber(usr.ID)

	if err := tx.Update(usr); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Update(topic); err != nil {
		return errors.WithStack(err)
	}

	c.Set("topic", topic)
	c.Set("category", topic.Category)
	c.Set("replies", &topic.Replies)
	return c.Redirect(302, "/topics/detail/%s", topic.ID)
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
	replies := make(models.Replies, 0, len(topic.Replies))
	for i := range topic.Replies {
		reply, err := loadReply(c, topic.Replies[i].ID.String())
		if err != nil {
			return nil, c.Error(404, err)
		}
		if reply.Deleted {
			continue
		}
		replies = append(replies, *reply)
	}
	topic.Replies = replies
	sort.Sort(topic.Replies)
	return topic, nil
}

func newReplyNotify(c buffalo.Context, topic *models.Topic, reply *models.Reply) error {
	set := make(map[uuid.UUID]struct{})
	for _, usr := range topic.Subscribers {
		set[usr] = struct{}{}
	}
	set[reply.AuthorID] = struct{}{}

	cat := new(models.Category)
	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Find(cat, topic.CategoryID); err != nil {
		return errors.WithStack(err)
	}
	for _, usr := range cat.Subscribers {
		set[usr] = struct{}{}
	}

	users := new(models.Users)
	if err := tx.All(users); err != nil {
		return errors.WithStack(err)
	}

	var recpts []models.User
	for _, usr := range *users {
		if _, ok := set[usr.ID]; !ok {
			continue
		}
		recpts = append(recpts, usr)
	}

	err := mailers.NewReplyNotify(c, topic, reply, recpts)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
