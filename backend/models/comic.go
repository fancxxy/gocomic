package models

import (
	"time"

	"github.com/fancxxy/gocomic/backend/rpc"
)

// Comic mapping table comics
type Comic struct {
	ID       int        `gorm:"primary_key" json:"id"`
	URL      string     `gorm:"size:128;not null;unique_index" json:"url"`
	Title    string     `gorm:"size:128;not null;unique_index:uix_comics_title_source" json:"title"`
	Source   string     `gorm:"size:32;not null;unique_index:uix_comics_title_source" json:"source"`
	Cover    string     `gorm:"size:128;not null" json:"cover"`
	Summary  string     `gorm:"type:text" json:"summary"`
	Latest   string     `json:"latest"`
	Created  Time       `gorm:"type:timestamp;not null;default:current_timestamp" json:"created"`
	Chapters []*Chapter `gorm:"foreignkey:ComicID;association_foreignkey:ID" json:"chapters"`
	Users    []*User    `gorm:"many2many:subscribers" json:"users,omitempty"`
}

// AllComics function get all comics information
func AllComics() []*Comic {
	var comics []*Comic
	db.Find(&comics)
	return comics
}

// QueryComic function query comic record from database
func QueryComic(source, title string, limit, page int) *Comic {
	var comic = new(Comic)
	if db.Where("title = ? and source = ?", title, source).First(comic).RecordNotFound() {
		return nil
	}

	var chapters []*Chapter
	paginate(limit, page).Where("comic_id = ?", comic.ID).Find(&chapters)
	comic.Chapters = chapters
	return comic
}

// GetComic function get comic record from database
func GetComic(id int, limit, page int) *Comic {
	comic := &Comic{ID: id}
	if db.First(comic).RecordNotFound() {
		return nil
	}

	var chapters []*Chapter
	paginate(limit, page).Where("comic_id = ?", comic.ID).Find(&chapters)
	comic.Chapters = chapters
	return comic
}

// AddComic function insert comic record and all chapter records
func AddComic(response *rpc.Comic, user *User) *Comic {
	comic := &Comic{
		URL:     response.URL,
		Title:   response.Title,
		Source:  response.Source,
		Cover:   response.Cover,
		Summary: response.Summary,
		Latest:  response.Latest,
	}

	var chapters []*Chapter
	for index, ctitle := range response.Indexes {
		chapters = append(chapters, &Chapter{
			URL:     response.Chapters[ctitle],
			Number:  index + 1,
			Title:   ctitle,
			Created: Time{time.Now()},
		})
	}
	comic.Chapters = chapters
	comic.Users = []*User{user}

	db.Create(comic)
	return comic
}

// Update insert new chapter information into chapters table and update latest column
func (c *Comic) Update(response *rpc.Update) bool {
	if len(response.Indexes) == 0 {
		c.Chapters = nil
		return false
	}

	var max int
	row := db.Table("chapters").Where("comic_id = ?", c.ID).Select("max(number)").Row()
	row.Scan(&max)
	var chapters []*Chapter
	for index, ctitle := range response.Indexes {
		chapters = append(chapters, &Chapter{
			URL:     response.Chapters[ctitle],
			Number:  max + index + 1,
			Title:   ctitle,
			Created: Time{time.Now()},
		})
	}
	c.Chapters = chapters
	c.Latest = response.Latest
	db.Save(c)
	return true
}
