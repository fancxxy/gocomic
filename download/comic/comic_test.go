package comic

import (
	"testing"
)

func TestComic(t *testing.T) {
	url, err := SearchComic("腾讯动漫", "航海王")
	if err != nil {
		t.Errorf("SearchComic(website, title) error: %v", err)
	}

	comic, err := NewComic(url)
	if err != nil {
		t.Errorf("NewComic(url) error: %v", err)
	}

	err = comic.DownloadCover(0.75)
	if err != nil {
		t.Errorf("Comic.DownloadCover(scale) error: %v", err)
	}

	// comic.Download()
	comic.Download(1, 10, 100, 1000)
}
