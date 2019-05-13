package parser

import (
	"net/url"
	"strings"
)

// Parser is the interface to wrap website methods
type Parser interface {
	Comic(url string) (map[string]interface{}, error)
	Chapter(url string) (map[string]interface{}, error)
	Download(chan map[string]string, chan string) (map[string]string, error)
	Search(string) (string, error)
	Less(string, string) bool
	Match(int, string) bool
	Label() string
}

const guard = 20

var instances = make(map[string]Parser)

func init() {
	tencent := newTencent()
	instances[tencent.Label()] = tencent
}

// GetParser find the Parser using input parameter, it can be Parser's name or domain
func GetParser(v string) Parser {
	var name, domain string
	u, err := url.Parse(v)
	if err != nil || u.Scheme == "" {
		name = v
	} else {
		domain = u.Host
	}

	for label, instance := range instances {
		splits := strings.Split(label, ": ")
		n, d := splits[0], splits[1]
		if name == n {
			return instance
		}
		if domain != "" && strings.Contains(d, domain) {
			return instance
		}
	}

	return nil
}

// Support return supported website lists
func Support() []string {
	var support []string
	for label := range instances {
		support = append(support, label)
	}
	return support
}
