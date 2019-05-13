package comic

import (
	"fmt"
	"path/filepath"

	"github.com/fancxxy/gocomic/download/parser"
)

// NewChapter create chapter instance
func NewChapter(url string) (*Chapter, error) {
	if url == "" {
		return nil, fmt.Errorf("url is mandatory")
	}

	p := parser.GetParser(url)
	if p == nil {
		return nil, fmt.Errorf("cannot find parser for url %s", url)
	}

	chapter := &Chapter{URL: url, home: home, parser: p}

	values, err := p.Chapter(url)
	if err != nil {
		return nil, err
	}

	chapter.initialize(values)
	return chapter, nil
}

// Chapter contains chapter informations
type Chapter struct {
	URL      string
	Title    string
	Ctitle   string
	pictures []string
	path     string
	home     string
	parser   parser.Parser
}

func (c *Chapter) initialize(values map[string]interface{}) {
	for field, value := range values {
		switch field {
		case "title":
			c.Title = value.(string)
		case "ctitle":
			c.Ctitle = value.(string)
		case "pictures":
			c.pictures = value.([]string)
		case "path":
			c.path = value.(string)
		}
	}
}

// Home reset chapter home directory
func (c *Chapter) Home(paths ...string) (string, error) {
	home, err := resolvePath(paths)
	if err != nil {
		return c.home, err
	}

	if home != "" {
		c.home = home
	}
	return c.home, nil
}

// Download download this chapter resource to disk
func (c *Chapter) Download() (map[string]string, error) {
	path := filepath.Join(c.home, c.path)
	err := createDir(path)
	if err != nil {
		return nil, err
	}

	queue := make(chan map[string]string, guard*4)
	go func() {
		for _, url := range c.pictures {
			picture := make(map[string]string)
			picture["url"] = url
			picture["curl"] = c.URL
			picture["filepath"] = path
			queue <- picture
		}
		close(queue)
	}()

	print := make(chan string)
	return c.parser.Download(queue, print)
}
