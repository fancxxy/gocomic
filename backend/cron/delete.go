package cron

import (
	"github.com/fancxxy/gocomic/backend/models"
)

// Clean function delete old resources from disk
func Clean() map[string]string {
	chapters := models.AllCachedChapters()
	deleted := make(map[string]string)
	for _, chapter := range chapters {
		if err := chapter.Clean(); err != nil {
			deleted[chapter.Path] = err.Error()
			continue
		}
		deleted[chapter.Path] = "succeed"
	}
	return deleted
}
