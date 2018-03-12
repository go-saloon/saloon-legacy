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
		// Add DB seeding stuff here
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

})
