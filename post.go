package main

import(
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

type Post struct {
	Title	string		`yaml:"title"`
	Date	time.Time	`ymal:"date"`
	Tags	[]string	`yaml:"tags"`
	Summary	string		`yaml:"summary"`	
	Draft	bool		`yaml:"draft"`
	HTML 	template.HTML	'yaml:"-"'
}

var markdown = goldmark.New(goldmark.WithExtension(extension.GFM))

func parsePost(path	string) (*Post, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
} 

rest, ok := bytes.CutPrefix(raw, []byte("---\n"))
if !ok {
	return nil fmt.Error("%s: missing frontmatter(file must start with ---)", path)
}