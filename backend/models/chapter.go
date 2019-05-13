package models

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fancxxy/gocomic/backend/rpc"
)

// Chapter mapping table chapters
type Chapter struct {
	ID       int        `gorm:"primary_key" json:"-"`
	URL      string     `gorm:"size:128;not null;unique_index" json:"url"`
	Number   int        `gorm:"not null;unique_index:uix_chapters_number_comic_id" json:"index"`
	Title    string     `gorm:"size:128;not null" json:"title"`
	Cached   bool       `gorm:"not null;default:false" json:"-"`
	Path     string     `gorm:"size:256"`
	ComicID  int        `gorm:"not null;unique_index:uix_chapters_number_comic_id" json:"-"`
	Created  Time       `gorm:"type:timestamp;not null;default:current_timestamp" json:"created"`
	Comic    *Comic     `gorm:"foreignkey:ComicID;association_foreignkey:ID" json:"-"`
	Pictures []string   `gorm:"-" json:"pictures"`
	PictureS []*Picture `gorm:"foreignkey:ChapterID;association_foreignkey:ID" json:"-"`
}

// AllCachedChapters function get all cached chapter information
func AllCachedChapters() []*Chapter {
	var chapters []*Chapter
	db.Where("cached = 1 and date_sub(curdate(), INTERVAL 7 DAY) >= date(`created`)").Find(&chapters)
	return chapters
}

// GetChapter function get comic record from database
func GetChapter(id, cid int) *Chapter {
	var chapter = new(Chapter)
	if db.Where("comic_id = ? and number = ?", id, cid).First(chapter).RecordNotFound() {
		return nil
	}

	if chapter.Cached {
		var pictures []*Picture
		db.Where("chapter_id = ?", chapter.ID).Order("number asc").Find(&pictures).Pluck("filename", &chapter.Pictures)
	}
	return chapter
}

// Save insert path of the resources into pictures table
func (c *Chapter) Save(response *rpc.Chapter) {
	if c.Path != "" {
		c.Cached = true
		c.Created = Time{time.Now()}
		db.Save(c)
		return
	}

	var pictures []*Picture

	var files []string
	for file := range response.Pictures {
		files = append(files, file)
	}
	sort.Strings(files)

	for index, file := range files {
		pictures = append(pictures, &Picture{
			Number:   index + 1,
			Filename: response.Pictures[file],
		})
		c.Pictures = append(c.Pictures, response.Pictures[file])
	}
	c.PictureS = pictures
	c.Path = filepath.Dir(c.Pictures[0])
	c.Cached = true
	db.Save(c)
}

// Clean function delete comic files and update cached column to false
func (c *Chapter) Clean() error {
	if err := os.RemoveAll(c.Path); err != nil {
		return err
	}
	c.Cached = false
	db.Save(c)
	return nil
}
