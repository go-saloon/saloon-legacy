// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/go-saloon/saloon/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

const indexName = "saloon.bleve.search"

var wrkr worker.Worker
var index bleve.Index

func init() {
	os.RemoveAll(indexName)
	var err error
	index, err = bleve.Open(indexName)
	if err == bleve.ErrorIndexPathDoesNotExist {
		index, err = bleve.New(indexName, bleve.NewIndexMapping())
		if err != nil {
			log.Fatalf("could not create bleve index: %v", err)
		}
	}

	wrkr = worker.NewSimple()
	wrkr.Register("index-db", func(args worker.Args) error {
		return indexDB()
	})

	go runIndex()
}

func runIndex() {
	tick := time.NewTicker(30 * time.Minute)
	defer tick.Stop()

	run := func() {
		wrkr.Perform(worker.Job{
			Queue:   "default",
			Handler: "index-db",
		})
	}

	run()
	for range tick.C {
		run()
	}
}

func indexDB() error {
	type indexedTopic struct {
		ID      uuid.UUID
		Title   string
		Content string
	}

	type indexedReply struct {
		ID      uuid.UUID
		Content string
	}

	return models.DB.Transaction(func(tx *pop.Connection) error {
		topics := new(models.Topics)
		if err := tx.All(topics); err != nil {
			return errors.WithStack(err)
		}
		for _, t := range *topics {
			err := index.Index("topics/detail/"+t.ID.String(), t)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		replies := new(models.Replies)
		if err := tx.All(replies); err != nil {
			return errors.WithStack(err)
		}
		for _, r := range *replies {
			err := index.Index(fmt.Sprintf("topics/detail/%s#%s", r.TopicID, r.ID), r)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}

func Search(c buffalo.Context) error {
	if c.Param("query") == "" {
		return c.Render(200, r.HTML("search"))
	}

	query := bleve.NewQueryStringQuery(c.Param("query"))
	req := bleve.NewSearchRequest(query)
	req.Size = 100
	req.Highlight = bleve.NewHighlight()

	res, err := index.Search(req)
	if err != nil {
		return errors.WithStack(err)
	}
	c.Set("results", res)
	return c.Render(200, r.HTML("search"))
}
