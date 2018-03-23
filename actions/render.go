// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"math"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
	"github.com/shurcooL/github_flavored_markdown"
)

var r *render.Engine
var assetsBox = packr.NewBox("../public/assets")

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",
		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),
		AssetsBox:    assetsBox,
		// Add template helpers here:
		Helpers: render.Helpers{
			"csrf": func() template.HTML {
				return template.HTML("<input name=\"authenticity_token\" value=\"<%= authenticity_token %>\" type=\"hidden\">")
			},
			"timeSince": timeSince,
			"score": func(f float64) string {
				return fmt.Sprintf("%.2f%%", f*100)
			},
			"markdown": markdownHelper,
		},
	})
}

func timeSince(created time.Time, ctx plush.HelperContext) string {
	if true && false {
		return created.UTC().Format(time.RFC3339)
	}

	now := time.Now().UTC()
	delta := now.Sub(created.UTC())
	days := int(math.Abs(delta.Hours()) / 24)
	if days > 30 {
		return created.Format("2006-02-01")
	}
	if days >= 1 {
		return fmt.Sprintf("%dj", days)
	}
	if delta.Hours() >= 1 {
		return fmt.Sprintf("%dh", int(delta.Hours()))
	}
	if delta.Minutes() >= 1 {
		return fmt.Sprintf("%dm", int(delta.Minutes()))
	}
	return fmt.Sprintf("%ds", int(delta.Seconds()))
}

func markdownHelper(body string, help plush.HelperContext) (template.HTML, error) {
	var err error
	if help.HasBlock() {
		body, err = help.Block()
		if err != nil {
			return "", err
		}
	}

	return markdown(body)
}

func markdown(body string) (template.HTML, error) {
	type segment struct {
		code bool
		data []byte
	}

	var (
		code = []byte("```")
		nln  = []byte("\n")
	)

	var blocks []segment
	r := bufio.NewScanner(bytes.NewReader([]byte(body)))
	for r.Scan() {
		txt := r.Bytes()
		switch {
		case bytes.HasPrefix(txt, code):
			switch n := len(blocks); n {
			case 0:
				// first code-block
				blk := segment{code: true}
				blk.data = append(blk.data, txt...)
				blk.data = append(blk.data, nln...)
				blocks = append(blocks, blk)
			default:
				blk := &blocks[n-1]
				switch blk.code {
				case true:
					// closing code-block
					blk.data = append(blk.data, txt...)
					blk.data = append(blk.data, nln...)
					blocks = append(blocks, segment{code: false})
				case false:
					// opening code-block
					blk := segment{code: true}
					blk.data = append(blk.data, txt...)
					blk.data = append(blk.data, nln...)
					blocks = append(blocks, blk)
				}
			}
		default:
			switch n := len(blocks); n {
			case 0:
				blk := segment{code: false}
				blk.data = append(blk.data, txt...)
				blk.data = append(blk.data, nln...)
				blocks = append(blocks, blk)
			default:
				blk := &blocks[n-1]
				blk.data = append(blk.data, txt...)
				blk.data = append(blk.data, nln...)
			}
		}
	}

	out := new(bytes.Buffer)
	for _, blk := range blocks {
		var data []byte
		switch blk.code {
		case true:
			data = blk.data
		case false:
			data = []byte(template.HTMLEscapeString(string(blk.data)))
		}
		_, err := out.Write(github_flavored_markdown.Markdown(data))
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	return template.HTML(out.Bytes()), nil
}
