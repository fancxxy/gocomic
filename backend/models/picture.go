package models

// Picture mapping table pictures
type Picture struct {
	ID        int64  `gorm:"primary_key"`
	Number    int    `gorm:"not null" json:"number"`
	Filename  string `gorm:"size:256"`
	ChapterID int    `gorm:"not null;index"`
	Created   Time   `gorm:"type:timestamp;not null;default:current_timestamp"`
}
