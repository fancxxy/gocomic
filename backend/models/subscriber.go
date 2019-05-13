package models

import (
	"fmt"
)

// Subscriber mapping table subscribers
type Subscriber struct {
	UserID  int  `gorm:"primary_key;index;auto_increment:false"`
	ComicID int  `gorm:"primary_key;index;auto_increment:false"`
	Created Time `gorm:"type:timestamp;not null;default:current_timestamp"`
	User    *User
	Comic   *Comic
}

// Subscribe function insert record into subscriber table
func Subscribe(user *User, comic *Comic) error {
	if Subscribed(user, comic) {
		return fmt.Errorf("already subscribe the comic")
	}
	db.Create(&Subscriber{User: user, Comic: comic})
	return nil
}

// Unsubscribe function delete record from subscriber table
func Unsubscribe(user *User, comic *Comic) error {
	if !Subscribed(user, comic) {
		return fmt.Errorf("not subscribe the comic yet")
	}
	db.Where("user_id = ? and comic_id = ?", user.ID, comic.ID).Delete(Subscriber{})
	return nil
}

// Subscribed function check the record exists or not
func Subscribed(user *User, comic *Comic) bool {
	var subscriber = new(Subscriber)
	return !db.Where("user_id = ? and comic_id = ?", user.ID, comic.ID).First(subscriber).RecordNotFound()
}

// SubscribeDate function return the Time field
func SubscribeDate(user *User, comic *Comic) Time {
	var subscriber = new(Subscriber)
	if db.Where("user_id = ? and comic_id = ?", user.ID, comic.ID).First(subscriber).RecordNotFound() {
		return Time{}
	}
	return subscriber.Created
}
