package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"go.yaml.in/yaml/v3"
)

// Post is one Markdown file from content/posts, parsed and rendered.
type Post struct {
	Slug    string        // from the filename: hello-world.md -> hello-world
	Title   string        `yaml:"title"`
	Date    time.Time     `yaml:"date"`
	Tags    []string      `yaml:"tags"`
	Summary string        `yaml:"summary"`
	Draft   bool          `yaml:"draft"`
	HTML    template.HTML `yaml:"-"` // the rendered body
}

// One shared renderer. GFM adds tables, strikethrough and autolinks.
var markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))

// parsePost reads one .md file, splits off the frontmatter, and renders the body.
func parsePost(path string) (*Post, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// The file must look like: ---\n<yaml>\n---\n<markdown>
	rest, ok := bytes.CutPrefix(raw, []byte("---\n"))
	if !ok {
		return nil, fmt.Errorf("%s: missing frontmatter (file must start with ---)", path)
	}
	meta, body, ok := bytes.Cut(rest, []byte("\n---\n"))
	if !ok {
		return nil, fmt.Errorf("%s: frontmatter is not closed with ---", path)
	}

	post := &Post{Slug: strings.TrimSuffix(filepath.Base(path), ".md")}
	if err := yaml.Unmarshal(meta, post); err != nil {
		return nil, fmt.Errorf("%s: bad frontmatter: %w", path, err)
	}
	if post.Title == "" || post.Date.IsZero() {
		return nil, fmt.Errorf("%s: title and date are required", path)
	}

	var buf bytes.Buffer
	if err := markdown.Convert(body, &buf); err != nil {
		return nil, fmt.Errorf("%s: render markdown: %w", path, err)
	}
	post.HTML = template.HTML(buf.String())
	return post, nil
}
