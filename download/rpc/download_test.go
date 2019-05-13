package rpc

import (
	"fmt"
	"net/rpc"
	"testing"
)

func TestDownloadComic(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:7700")
	if err != nil {
		t.Errorf("tpc.DialHTTP error: %v", err)
	}

	request := &Parameter{
		Website: "腾讯漫画",
		Comic:   "航海王",
	}

	var response Comic
	err = client.Call("Download.Comic", request, &response)
	if err != nil {
		t.Errorf("client.Call error: %v", err)
	}

	fmt.Printf("Download.Comic response: %v", response)
}

func TestDownloadChapter(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:7700")
	if err != nil {
		t.Errorf("tpc.DialHTTP error: %v", err)
	}

	request := &Parameter{
		Chapter: "https://ac.qq.com/ComicView/index/id/505430/cid/957",
		Path:    "~",
	}

	var response Chapter
	err = client.Call("Download.Chapter", request, &response)
	if err != nil {
		t.Errorf("client.Call error: %v", err)
	}

	fmt.Printf("Download.Chapter response: %v", response)
}

func TestDownloadUpdate(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:7700")
	if err != nil {
		t.Errorf("tpc.DialHTTP error: %v", err)
	}

	request := &Parameter{
		Comic:  "https://ac.qq.com/Comic/comicInfo/id/505430",
		Latest: "航海王：第937话 强盗桥的牛鬼丸",
	}

	var response Update
	err = client.Call("Download.Update", request, &response)
	if err != nil {
		t.Errorf("client.Call error: %v", err)
	}

	fmt.Printf("Download.Update response: %v", response)
}
