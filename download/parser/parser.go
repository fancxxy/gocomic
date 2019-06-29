package parser

import (
	"net/url"
	"strings"
)

// Parser is the interface to wrap website methods
type Parser interface {
	// Chapter parse chapter page
	Chapter(url string) (map[string]interface{}, error)
	// Comic parse comic main page return useful information
	Comic(url string) (map[string]interface{}, error)
	// Filename get picture name from url
	Filename(string) string
	// Label return name + domain
	Label() string
	// Less compare two chapters' title
	Less(string, string) bool
	// Match check chapter title match input index or not
	Match(int, string) bool
	// Search search title in website and return the url of matched comic
	Search(string) (string, error)
}

var instances = make(map[string]Parser)

func init() {
	tencent := newTencent()
	// Label is 腾讯动漫: https://ac.qq.com/
	instances[tencent.Label()] = tencent
}

// GetParser find the Parser using input parameter, parameter can be Parser's name or domain
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
