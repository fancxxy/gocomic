package comic

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/fancxxy/gocomic/download/network"
	"github.com/fancxxy/gocomic/download/parser"
)

var (
	home string
)

const (
	guard = 5
	scale = 0.75
)

func init() {
	usr, err := user.Current()
	if err == nil {
		home = filepath.Join(usr.HomeDir, "Comics")
	}
}

// SearchComic function return the real url about the comic from the website
func SearchComic(website, title string) (string, error) {
	if website == "" || title == "" {
		return "", fmt.Errorf("website and title is mandatory")
	}

	p := parser.GetParser(website)
	if p == nil {
		return "", fmt.Errorf("cannot find parser for website %s", website)
	}

	u, err := p.Search(title)
	if err != nil {
		return "", err
	}

	return u, nil
}

// NewComic create comic instance
func NewComic(url string) (*Comic, error) {
	if url == "" {
		return nil, fmt.Errorf("url is mandatory")
	}

	p := parser.GetParser(url)
	if p == nil {
		return nil, fmt.Errorf("cannot find parser for url %s", url)
	}

	comic := &Comic{URL: url, home: home, parser: p}

	values, err := p.Comic(url)
	if err != nil {
		return nil, err
	}

	comic.initialize(values)
	comic.sort()

	if len(comic.Ctitles) < 1 {
		return nil, fmt.Errorf("chapter list is empty")
	}
	comic.Latest = comic.Ctitles[len(comic.Ctitles)-1]

	return comic, nil
}

// Comic contains comic informations
type Comic struct {
	URL      string
	Title    string
	Chapters map[string]string
	Summary  string
	Cover    string
	Ctitles  []string
	Latest   string
	Date     string
	path     string
	home     string
	parser   parser.Parser
}

func (c *Comic) initialize(values map[string]interface{}) {
	for field, value := range values {
		switch field {
		case "title":
			c.Title = value.(string)
		case "chapters":
			c.Chapters = value.(map[string]string)
			for title := range c.Chapters {
				c.Ctitles = append(c.Ctitles, title)
			}
		case "summary":
			c.Summary = value.(string)
		case "cover":
			c.Cover = value.(string)
		case "path":
			c.path = value.(string)
		case "date":
			c.Date = value.(string)
		}
	}
}

// Home reset comic home directory
func (c *Comic) Home(paths ...string) (string, error) {
	home, err := resolvePath(paths)
	if err != nil {
		return c.home, err
	}

	if home != "" {
		c.home = home
	}
	return c.home, nil
}

func (c *Comic) curls(nums ...int) []string {
	var urls []string
	var titles []string
	if len(nums) == 0 {
		for _, title := range c.Ctitles {
			urls = append(urls, c.Chapters[title])
		}
		return urls
	}

	sort.Ints(nums)

	var i, length = 0, len(c.Ctitles)
	for _, num := range nums {
		for {
			if i >= length {
				break
			}
			title := c.Ctitles[i]
			if c.parser.Match(num, title) {
				titles = append(titles, title)
				break
			}
			i++
		}
	}

	sortTitle(titles, c.parser)
	for _, title := range titles {
		urls = append(urls, c.Chapters[title])
	}
	return urls
}

func (c *Comic) sort() {
	sortTitle(c.Ctitles, c.parser)
}

// DownloadCover download cover image to disk
func (c *Comic) DownloadCover(s ...float64) error {
	var tempScale float64
	switch len(s) {
	case 0:
		tempScale = scale
	case 1:
		tempScale = s[0]
	case 2:
		return fmt.Errorf("too many paramters")
	}

	if c.Cover == "" {
		return fmt.Errorf("cannot get cover")
	}

	path := filepath.Join(c.home, c.path)
	err := createDir(path)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(path, "cover.jpg"))
	if err != nil {
		return err
	}
	defer file.Close()

	request := network.New()
	response, err := request.Get(c.Cover)
	if err != nil {
		return err
	}

	err = clipCover(response.ToReader(), file, tempScale)
	if err != nil {
		return err
	}

	return nil
}

// DownloadInCmd download resource to disk
func (c *Comic) DownloadInCmd(nums ...int) {
	var wg sync.WaitGroup

	urls := c.curls(nums...)
	if len(urls) == 0 {
		return
	}

	urlQueue := make(chan string, guard)
	pictureQueue := make(chan map[string]string, guard*4)
	printQueue := make(chan string, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		download(pictureQueue, printQueue, guard*4, c.parser.Filename)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			message, ok := <-printQueue
			if !ok {
				return
			}
			fmt.Println(message)
		}
	}()

	var wait sync.WaitGroup
	wait.Add(guard)
	for i := 0; i < guard; i++ {
		go func() {
			defer wait.Done()
			for {
				url, ok := <-urlQueue
				if !ok {
					return
				}

				chapter, err := NewChapter(url)
				if err != nil {
					continue
				}

				path := filepath.Join(chapter.home, chapter.path)
				err = createDir(path)
				if err != nil {
					continue
				}

				picture := make(map[string]string)
				picture["print"] = ""
				picture["title"] = chapter.Title
				picture["ctitle"] = chapter.Ctitle
				picture["path"] = path
				pictureQueue <- picture

				for _, url := range chapter.pictures {
					picture := make(map[string]string)
					picture["url"] = url
					picture["curl"] = chapter.URL
					picture["filepath"] = path
					pictureQueue <- picture
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		wait.Wait()
		close(pictureQueue)
	}()

	go func() {
		for _, url := range urls {
			urlQueue <- url
		}
		close(urlQueue)
	}()

	wg.Wait()
}

// Download download resource to disk, return all errors
func (c *Comic) Download(nums ...int) map[string]error {
	errs := make(map[string]error)

	urls := c.curls(nums...)
	if len(urls) == 0 {
		return errs
	}

	type result struct {
		url string
		err error
	}

	urlTitleMap := make(map[string]string)
	for title, url := range c.Chapters {
		urlTitleMap[url] = title
	}

	urlQueue := make(chan string, guard)
	resultQueue := make(chan *result, guard)

	for i := 0; i < guard; i++ {
		go func() {
			for {
				url, ok := <-urlQueue
				if !ok {
					return
				}
				chapter, err := NewChapter(url)
				if err != nil {
					resultQueue <- &result{url, err}
				} else {
					chapter.home = c.home
					_, err = chapter.Download()
					resultQueue <- &result{url, err}
				}
			}
		}()
	}

	go func() {
		for _, url := range urls {
			urlQueue <- url
		}
		close(urlQueue)
	}()

	for i := 0; i < len(urls); i++ {
		result, ok := <-resultQueue
		if ok && result.err != nil {
			errs[urlTitleMap[result.url]] = result.err
		}
	}

	return errs
}

func download(queue chan map[string]string, print chan string, guard int, filename func(string) string) (map[string]string, error) {
	var (
		wg      sync.WaitGroup
		syncMap sync.Map
		errors  []string
		request = network.New()
	)

	wg.Add(guard)
	for i := 0; i < guard; i++ {
		go func() {
			defer wg.Done()
			for {
				picture, ok := <-queue
				if !ok {
					return
				}

				if _, ok := picture["print"]; ok {
					print <- fmt.Sprintf("  %s: %s ==> %s", picture["title"], picture["ctitle"], picture["path"])
					continue
				}

				url := picture["url"]
				referer := picture["curl"]
				response, err := request.Get(url, network.Header{
					"Referer": referer,
				})
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s", url))
					continue
				}

				name := filename(url)
				wholename := filepath.Join(picture["filepath"], name)
				err = response.ToFile(wholename)
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s", url))
					continue
				}
				syncMap.Store(name, wholename)
			}
		}()
	}

	wg.Wait()
	close(print)

	var resources = make(map[string]string)
	syncMap.Range(func(k, v interface{}) bool {
		resources[k.(string)] = v.(string)
		return true
	})

	if len(errors) != 0 {
		return resources, fmt.Errorf("missing: [%s]", strings.Join(errors, ", "))
	}
	return resources, nil
}
