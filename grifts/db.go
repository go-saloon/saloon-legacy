package grifts

import (
	"math/rand"

	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/pop"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		for _, name := range []string{
			"db:setup",
			"db:seed:create-forum",
			"db:seed:create-categories",
		} {
			err := grift.Run(name, c)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	grift.Desc("seed:create-categories", "Create a few default categories")
	grift.Add("seed:create-categories", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			for _, cat := range []*models.Category{
				{Title: "Category-1"},
				{Title: "Category-2"},
			} {
				err := tx.Create(cat)
				if err != nil {
					return errors.WithStack(err)
				}
			}
			users := new(models.Users)
			if err := tx.All(users); err != nil {
				return errors.WithStack(err)
			}
			cats := new(models.Categories)
			if err := tx.All(cats); err != nil {
				return errors.WithStack(err)
			}

			for _, usr := range *users {
				if usr.Admin {
					continue
				}
				i := rand.Intn(len(*cats))
				cat := (*cats)[i]
				usr.Subscriptions = append(usr.Subscriptions, cat.ID)
				if err := tx.Update(&usr); err != nil {
					return errors.WithStack(err)
				}
				cat.Subscribers = append(cat.Subscribers, usr.ID)
				if err := tx.Update(&cat); err != nil {
					return errors.WithStack(err)
				}
			}
			return nil
		})
	})

	grift.Desc("seed:create-forum", "Create forum welcome message")
	grift.Add("seed:create-forum", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			forum := &models.Forum{
				BaseAddr:    "127.0.0.1:3000",
				Title:       "Alternatiba-63",
				Description: "a Forum for Alternatiba 63",
			}
			err := tx.Create(forum)
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
	})
})
