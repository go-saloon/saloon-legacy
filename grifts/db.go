package grifts

import (
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
