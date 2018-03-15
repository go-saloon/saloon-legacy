package grifts

import (
	"bytes"
	"flag"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

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
		fset := flag.NewFlagSet("create-forum", flag.ExitOnError)
		fname := fset.String("logo", "", "path to a logo image")
		title := fset.String("title", "Saloon", "title of the forum")
		descr := fset.String("descr", "A nice forum", "description of the forum")

		err := fset.Parse(c.Args)
		if err != nil {
			return errors.WithStack(err)
		}

		return models.DB.Transaction(func(tx *pop.Connection) error {
			var (
				err  error
				logo []byte
			)

			if *fname != "" {
				f, err := os.Open(*fname)
				if err != nil {
					return errors.WithStack(err)
				}
				defer f.Close()
				img, _, err := image.Decode(f)
				if err != nil {
					return errors.WithStack(err)
				}
				buf := new(bytes.Buffer)
				err = png.Encode(buf, img)
				if err != nil {
					return errors.WithStack(err)
				}
				logo = buf.Bytes()
			}
			forum := &models.Forum{
				Title:       *title,
				Description: *descr,
				Logo:        logo,
			}
			err = tx.Create(forum)
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
	})
})
