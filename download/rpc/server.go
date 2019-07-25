package rpc

import (
	"fmt"
	"io"
	"net"

	"github.com/fancxxy/gocomic/download/comic"
	"github.com/fancxxy/gocomic/download/rpc/download"
	"google.golang.org/grpc"
)

const (
	port = ":7070"
)

type server struct{}

func (s *server) Comic(stream download.Download_ComicServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		url, err := comic.SearchComic(request.Website, request.Comic)
		if err != nil {
			return err
		}

		comic, err := comic.NewComic(url)
		if err != nil {
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
			return err
		}

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

		response := &download.ChapterResponse{
			Url:      chapter.URL,
			Title:    chapter.Title,
			Ctitle:   chapter.Ctitle,
			Pictures: pictures,
		}
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
			return err
		}

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

		response := &download.UpdateResponse{
			Chapters: chapters,
			Indexes:  indexes,
		}
		stream.Send(response)
	}
}

// Start rpc server
func Start() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	fmt.Println("listen and serve:", port)

	g := grpc.NewServer()
	download.RegisterDownloadServer(g, &server{})
	g.Serve(lis)

	return nil
}
