package main

import (
	"strings"
	"testing"
)

func TestParsePost(t *testing.T) {
	post, err := parsePost("content/posts/hello-world.md")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", post)

	if post.Slug != "hello-world" {
		t.Errorf("slug = %q, want hello-world", post.Slug)
	}
	if post.Date.Format("2006-01-02") != "2026-11-23" {
		t.Errorf("date = %v, want 2026-11-23", post.Date)
	}
	for _, want := range []string{"<h2", "<strong>bold</strong>", "<pre><code"} {
		if !strings.Contains(string(post.HTML), want) {
			t.Errorf("rendered HTML is missing %s", want)
		}
	}
}
