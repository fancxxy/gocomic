package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/fancxxy/gocomic/download/network"
	"github.com/robertkrimen/otto"
)

var javascript = `
	function Base() {
		_keyStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";
		this.decode = function (c) {
			var a = "", b, d, h, f, g, e = 0;

			for (c = c.replace(/[^A-Za-z0-9\+\/\=]/g, "");
				e < c.length;
			)b = _keyStr.indexOf(c.charAt(e++)), d = _keyStr.indexOf(c.charAt(e++)), f = _keyStr.indexOf(c.charAt(e++)), g = _keyStr.indexOf(c.charAt(e++)), b = b << 2 | d >> 4, d = (d & 15) << 4 | f >> 2, h = (f & 3) << 6 | g, a += String.fromCharCode(b), 64 != f && (a += String.fromCharCode(d)), 64 != g && (a += String.fromCharCode(h));
			return a = _utf8_decode(a)
		};
		_utf8_decode = function (c) {
			for (var a = "", b = 0, d = c1 = c2 = 0;
				b < c.length;
			)d = c.charCodeAt(b), 128 > d ? (a += String.fromCharCode(d), b++) : 191 < d && 224 > d ? (c2 = c.charCodeAt(b + 1), a += String.fromCharCode((d & 31) << 6 | c2 & 63), b += 2) : (c2 = c.charCodeAt(b + 1), c3 = c.charCodeAt(b + 2), a += String.fromCharCode((d & 15) << 12 | (c2 & 63) << 6 | c3 & 63), b += 3);
			return a
		}
	}

	function get(v) {
		return true
	}

	var window = new Object(); 
	window.Array = 1

	var document = new Object();
	document.children = 1
	document.getElementsByTagName = get
	document.getElementsById = get
	document.getElementssByName = get

	%s 
	%s

	var B = new Base(), T = DATA.split(''), N = window['nonce'], len, locate, str;
	N = N.match(/\d+[a-zA-Z]+/g);
	len = N.length;
	while (len--) {
		locate = parseInt(N[len]) & 255;
		str = N[len].replace(/\d+/g, '');
		T.splice(locate, str.length)
	}
	T = T.join('');
	_v = B.decode(T);
`

type pictures struct {
	Comic struct {
		Title string `json:"title"`
	} `json:"comic"`
	Chapter struct {
		Ctitle string `json:"cTitle"`
	} `json:"chapter"`
	Picture []struct {
		URL string `json:"url"`
	} `json:"picture"`
}

type tencent struct {
	name        string
	host        string
	search      string
	cssSelector map[string]string
	regex       map[string]*regexp.Regexp
	request     *network.Request
}

func newTencent() *tencent {
	t := new(tencent)

	t.name = "腾讯动漫"
	t.host = "https://ac.qq.com"
	t.search = "https://ac.qq.com/Comic/searchList/search/"
	t.cssSelector = map[string]string{
		"title":    ".works-intro-title > strong:nth-child(1)",
		"chapters": "#chapter",
		"chapter":  ".works-chapter-item",
		"nonce":    "body > script:nth-child(15)",
		"data":     "body > script:nth-child(30)",
		"summary":  ".works-intro-short",
		"cover":    ".works-cover > a:nth-child(1)",
		"search":   ".mod_book_name",
	}
	urlRegex, _ := regexp.Compile(`ac.qq.com/Comic/comicInfo/id/\d+`)
	curlRegex, _ := regexp.Compile(`ac.qq.com/ComicView/index/id/\d+/cid/\d+`)
	chapterRegex, _ := regexp.Compile(`\d+`)
	t.regex = map[string]*regexp.Regexp{
		"url":     urlRegex,
		"curl":    curlRegex,
		"chapter": chapterRegex,
	}
	t.request = network.New()

	return t
}

func (t *tencent) Label() string {
	return t.name + ": " + t.host + "/"
}

func (t *tencent) Comic(url string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	response, err := t.request.Get(url)
	if err != nil {
		return results, err
	}

	doc, err := goquery.NewDocumentFromReader(response.ToReader())
	if err != nil {
		return results, err
	}

	chapters := make(map[string]string)
	title := doc.Find(t.cssSelector["title"]).Text()
	doc.Find(t.cssSelector["chapters"]).Find(t.cssSelector["chapter"]).Each(func(_ int, s *goquery.Selection) {
		title, _ := s.Find("a").Attr("title")
		url, _ := s.Find("a").Attr("href")
		url = t.host + url
		chapters[title] = url
	})
	cover, _ := doc.Find(t.cssSelector["cover"]).Find("img").Attr("src")
	summary := strings.Trim(doc.Find(t.cssSelector["summary"]).Text(), "\n ")
	date := doc.Find(t.cssSelector["update"]).Find(".ui-pl10").Text()

	results["title"] = title
	results["cover"] = cover
	results["chapters"] = chapters
	results["summary"] = summary
	results["path"] = filepath.Join(t.name, title)
	results["date"] = date

	return results, nil
}

func (t *tencent) Chapter(url string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	response, err := t.request.Get(url)
	if err != nil {
		return results, err
	}

	doc, err := goquery.NewDocumentFromReader(response.ToReader())
	if err != nil {
		return results, err
	}

	nonce := doc.Find(t.cssSelector["nonce"]).Text()
	data := doc.Find(t.cssSelector["data"]).Text()

	if nonce == "" || data == "" {
		return results, fmt.Errorf("nonce or data is empty")
	}

	vm := otto.New()
	value, err := vm.Run(fmt.Sprintf(javascript, nonce, data))
	if err != nil {
		return results, err
	}

	var pics pictures
	err = json.Unmarshal([]byte(value.String()), &pics)

	var urls []string
	for _, picture := range pics.Picture {
		urls = append(urls, picture.URL)
	}

	results["title"] = pics.Comic.Title
	results["ctitle"] = pics.Chapter.Ctitle
	results["pictures"] = urls
	results["path"] = filepath.Join(t.name, pics.Comic.Title, pics.Chapter.Ctitle)

	return results, nil
}

func (t *tencent) Match(num int, title string) bool {
	matched, _ := regexp.MatchString(fmt.Sprintf("第%d话? ", num), title)
	return matched
}

func (t *tencent) Less(title1, title2 string) bool {
	num1, num2 := t.getChapterNum(title1), t.getChapterNum(title2)
	return num1 < num2
}

func (t *tencent) Download(queue chan map[string]string, print chan string) (map[string]string, error) {
	var wg sync.WaitGroup
	var errors []string

	var syncMap sync.Map
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
					print <- fmt.Sprintf("  %s: %-20s ==> %s", picture["title"], picture["ctitle"], picture["path"])
					continue
				}

				url := picture["url"]
				referer := picture["curl"]
				response, err := t.request.Get(url, network.Header{
					"Referer": referer,
				})
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s", url))
				}

				filename := url[strings.LastIndex(url, "_")+1 : len(url)-2]
				wholename := filepath.Join(picture["filepath"], filename)
				err = response.ToFile(wholename)
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s", url))
				}
				syncMap.Store(filename, wholename)
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

func (t *tencent) Search(title string) (string, error) {
	query := t.search + url.QueryEscape(title)
	var URL string

	request := network.New()
	response, err := request.Get(query)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(response.ToReader())
	doc.Find(t.cssSelector["search"]).Each(func(_ int, s *goquery.Selection) {
		T, _ := s.Find("a").Attr("title")
		U, _ := s.Find("a").Attr("href")
		if title == T {
			URL = t.host + U
			return
		}

	})

	if URL == "" {
		return "", fmt.Errorf("cannot find comic %s", title)
	}
	return URL, nil
}

func (t *tencent) getChapterNum(title string) int {
	s := t.regex["chapter"].FindString(title)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
