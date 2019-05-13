package cron

import (
	"github.com/fancxxy/gocomic/backend/models"
	"github.com/fancxxy/gocomic/backend/rpc"
)

// Update function update all comics
func Update() map[string]map[string]string {
	comics := models.AllComics()

	updates := make(map[string]map[string]string)
	for _, comic := range comics {
		response, err := rpc.UpdateComic(comic.URL, comic.Latest)
		if err != nil {
			website, ok := updates[comic.Source]
			if !ok {
				website = make(map[string]string)
			}
			website[comic.Title] = "failed"
			updates[comic.Source] = website
			continue
		}

		if !comic.Update(response) {
			website, ok := updates[comic.Source]
			if !ok {
				website = make(map[string]string)
			}
			website[comic.Title] = "-"
			updates[comic.Source] = website
			continue
		}

		website, ok := updates[comic.Source]
		if !ok {
			website = make(map[string]string)
		}
		website[comic.Title] = "succeed"
		updates[comic.Source] = website
	}

	return updates
}
