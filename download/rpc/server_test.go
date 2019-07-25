package rpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/fancxxy/gocomic/download/rpc/download"
	"google.golang.org/grpc"
)

func TestComic(t *testing.T) {
	addr := "127.0.0.1:7070"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connect server error: %v", err)
	}
	defer conn.Close()

	client := download.NewDownloadClient(conn)
	stream, err := client.Comic(context.Background())
	if err != nil {
		t.Errorf("receive stream error: %v", err)
	}

	var i int32
	for i = 0; i < 1; {
		request := &download.ComicRequest{
			Website: "腾讯动漫",
			Comic:   "航海王",
		}
		stream.Send(request)
		response, err := stream.Recv()
		if err != nil {
			t.Errorf("resp error: %v", err)
		}

		fmt.Println(response.Url)
		fmt.Println(response.Title)
		fmt.Println(response.Summary)
		fmt.Println(response.Cover)
		fmt.Println(response.Chapters)
		fmt.Println(response.Indexes)
		fmt.Println(response.Latest)
		fmt.Println(response.Source)
		i++
	}
}

func TestChapter(t *testing.T) {
	addr := "127.0.0.1:7070"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connect server error: %v", err)
	}
	defer conn.Close()

	client := download.NewDownloadClient(conn)
	stream, err := client.Chapter(context.Background())
	if err != nil {
		t.Errorf("receive stream error: %v", err)
	}

	var i int32
	for i = 0; i < 1; {
		request := &download.ChapterRequest{
			Chapter: "https://ac.qq.com/ComicView/index/id/505430/cid/966",
			Path:    "",
		}
		stream.Send(request)
		response, err := stream.Recv()
		if err != nil {
			t.Errorf("resp error: %v", err)
		}

		fmt.Println(response.Url)
		fmt.Println(response.Title)
		fmt.Println(response.Ctitle)
		fmt.Println(response.Pictures)
		i++
	}
}

func TestUpdate(t *testing.T) {
	addr := "127.0.0.1:7070"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Errorf("connect server error: %v", err)
	}
	defer conn.Close()

	client := download.NewDownloadClient(conn)
	stream, err := client.Update(context.Background())
	if err != nil {
		t.Errorf("receive stream error: %v", err)
	}

	var i int32
	for i = 0; i < 1; {
		request := &download.UpdateRequest{
			Comic:  "https://ac.qq.com/Comic/comicInfo/id/505430",
			Latest: "航海王：第947话 昆因的赌注",
		}
		stream.Send(request)
		response, err := stream.Recv()
		if err != nil {
			t.Errorf("resp error: %v", err)
		}

		fmt.Println(response.Chapters)
		fmt.Println(response.Indexes)
		i++
	}
}
