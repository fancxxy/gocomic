package rpc

import (
	"net/rpc"

	"github.com/fancxxy/gocomic/backend/config"
)

// Client is a point to rpc.Client
var Client *rpc.Client

// Init function initialize rpc client
func Init() error {
	address := config.Config.Download.Address
	var err error
	Client, err = rpc.DialHTTP("tcp", address)
	if err != nil {
		return err
	}

	return nil
}

// DownloadComic function remote call Download.Comic
func DownloadComic(website, comic string) (*Comic, error) {
	request := &Parameter{
		Website: website,
		Comic:   comic,
	}

	var response = new(Comic)
	err := Client.Call("Download.Comic", request, response)
	return response, err
}

// DownloadChapter function remote call Download.Chapter
func DownloadChapter(url, path string) (*Chapter, error) {
	request := &Parameter{
		Chapter: url,
		Path:    path,
	}

	var response = new(Chapter)
	err := Client.Call("Download.Chapter", request, response)
	return response, err
}

// UpdateComic function remote call Download.Update
func UpdateComic(url string, latest string) (*Update, error) {
	request := &Parameter{
		Comic:  url,
		Latest: latest,
	}

	var response = new(Update)
	err := Client.Call("Download.Update", request, &response)
	return response, err
}
