package rpc

import (
	"fmt"
	"io"
	"net"

	"github.com/fancxxy/gocomic/download/comic"
	"github.com/fancxxy/gocomic/download/rpc/download"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	port = "0.0.0.0:7070"
)

type server struct{}

func (s *server) Comic(stream download.Download_ComicServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Errorf("[comic] recv comic request failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		url, err := comic.SearchComic(request.Website, request.Comic)
		if err != nil {
			log.Errorf("[comic] search comic url failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		comic, err := comic.NewComic(url)
		if err != nil {
			log.Errorf("[comic] parse comic info failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		response := &download.ComicResponse{
			Url:      comic.URL,
			Title:    comic.Title,
			Cover:    comic.Cover,
			Summary:  comic.Summary,
			Chapters: comic.Chapters,
			Latest:   comic.Latest,
			Source:   request.Website,
			Indexes:  comic.Ctitles,
		}

		log.Debugf("[comic] call rpc succeed, request: %s, response: %s", oct2utf8(request.String()), oct2utf8(response.String()))
		stream.Send(response)
	}
}

func (s *server) Chapter(stream download.Download_ChapterServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Errorf("[chapter] recv chapter request failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		chapter, err := comic.NewChapter(request.Chapter)
		if err != nil {
			log.Errorf("[chapter] parse chapter info failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		if request.Path != "" {
			chapter.Home(request.Path)
		}

		pictures, err := chapter.Download()
		if err != nil {
			log.Errorf("[chapter] download chapter resource failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		response := &download.ChapterResponse{
			Url:      chapter.URL,
			Title:    chapter.Title,
			Ctitle:   chapter.Ctitle,
			Pictures: pictures,
		}

		log.Debugf("[chapter] call rpc succeed, request: %s, response: %s", oct2utf8(request.String()), oct2utf8(response.String()))
		stream.Send(response)
	}
}

func (s *server) Update(stream download.Download_UpdateServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Errorf("[update] recv update request failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		comic, err := comic.NewComic(request.Comic)
		if err != nil {
			log.Errorf("[update] parse comic info failed: %v, request: %s", err, oct2utf8(request.String()))
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
			err := fmt.Errorf("cannot find the latest chapter %v", request.Latest)
			log.Errorf("[update] match chapter title failed: %v, request: %s", err, oct2utf8(request.String()))
			return err
		}

		var chapters = make(map[string]string)
		var indexes []string
		for i := index; i < len(comic.Ctitles); i++ {
			title := comic.Ctitles[i]
			url := comic.Chapters[title]
			chapters[title] = url
			indexes = append(indexes, title)
		}

		response := &download.UpdateResponse{
			Chapters: chapters,
			Indexes:  indexes,
		}

		log.Debugf("[update] call rpc succeed, request: %s, response: %s", oct2utf8(request.String()), oct2utf8(response.String()))
		stream.Send(response)
	}
}

// Start rpc server
func Start(debug bool) error {
	initialzie(debug)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	log.Infof("gocomic run as rpc mode, listen and serve: %s", port)

	g := grpc.NewServer()
	download.RegisterDownloadServer(g, &server{})
	g.Serve(lis)

	return nil
}
