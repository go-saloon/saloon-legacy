// Copyright 2018 The go-saloon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import "testing"

func TestMarkdown(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
		want []byte
	}{
		{
			name: "empty",
			data: []byte(""),
			want: []byte(""),
		},
		{
			name: "with-div",
			data: []byte("<div>foo</div>\n\n<div>"),
			want: []byte("<p>&lt;div&gt;foo&lt;/div&gt;</p>\n\n<p>&lt;div&gt;</p>\n"),
		},
		{
			name: "with-code-block",
			data: []byte("<div>foo</div>\n\n```\nfunc() { \"hello\" }\n```\n<div>"),
			want: []byte("<p>&lt;div&gt;foo&lt;/div&gt;</p>\n<pre><code>func() { &#34;hello&#34; }\n</code></pre>\n<p>&lt;div&gt;</p>\n"),
		},
		{
			name: "with-code-block-at-start",
			data: []byte("```\nfunc() { \"hello\" }\n```\n<div>&mtimes;</div>"),
			want: []byte("<pre><code>func() { &#34;hello&#34; }\n</code></pre>\n<p>&lt;div&gt;&amp;mtimes;&lt;/div&gt;</p>\n"),
		},
		{
			name: "with-code-block-unmatched",
			data: []byte("foo\n\n```\n\nbar\n"),
			want: []byte("<p>foo</p>\n<p>```</p>\n\n<p>bar</p>\n"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			txt, err := markdown(string(tc.data))
			if err != nil {
				t.Fatal(err)
			}

			if string(txt) != string(tc.want) {
				t.Fatalf("error %q:\ngot: %q\nwant:%q\n", tc.name, txt, string(tc.want))
			}
		})
	}
}
