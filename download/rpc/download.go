package rpc

import (
	"fmt"

	"github.com/fancxxy/gocomic/download/comic"
)

// Comic get comic information and set response
func (d *Download) Comic(request *Parameter, response *Comic) error {
	url, err := comic.SearchComic(request.Website, request.Comic)
	if err != nil {
		return err
	}

	comic, err := comic.NewComic(url)
	if err != nil {
		return err
	}

	response.URL = comic.URL
	response.Title = comic.Title
	response.Cover = comic.Cover
	response.Summary = comic.Summary
	response.Chapters = comic.Chapters
	response.Latest = comic.Latest
	response.Source = request.Website
	response.Indexes = comic.Ctitles

	return nil
}

// Chapter get chapter information and set response
func (d *Download) Chapter(request *Parameter, response *Chapter) error {
	chapter, err := comic.NewChapter(request.Chapter)
	if err != nil {
		return err
	}

	if request.Path != "" {
		chapter.Home(request.Path)
	}

	pictures, err := chapter.Download()
	if err != nil {
		return err
	}

	response.URL = chapter.URL
	response.Title = chapter.Title
	response.Ctitle = chapter.Ctitle
	response.Pictures = pictures

	return nil
}

// Update return updated information about the comic
func (d *Download) Update(request *Parameter, response *Update) error {
	comic, err := comic.NewComic(request.Comic)
	if err != nil {
		return err
	}

	var index int
	for i, title := range comic.Ctitles {
		if title == request.Latest {
			index = i + 1
			break
		}
	}
	if index == 0 {
		return fmt.Errorf("cannot find the latest chapter %v", request.Latest)
	}

	var chapters = make(map[string]string)
	var indexes []string
	for i := index; i < len(comic.Ctitles); i++ {
		title := comic.Ctitles[i]
		url := comic.Chapters[title]
		chapters[title] = url
		indexes = append(indexes, title)
	}
	response.Latest = comic.Latest
	response.Chapters = chapters
	response.Indexes = indexes

	return nil
}
