package models

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User mapping table users
type User struct {
	ID       int      `gorm:"primary_key"`
	Email    string   `gorm:"size:64;not null;unique_index"`
	Username string   `gorm:"size:64;not null;unique_index"`
	Password string   `gorm:"size:128;not null"`
	Avatar   string   `gorm:"size:128"`
	Admin    bool     `gorm:"not null;default:false"`
	Blocked  bool     `gorm:"not null;default:false"`
	Created  Time     `gorm:"type:timestamp;not null;default:current_timestamp"`
	Comics   []*Comic `gorm:"many2many:subscribers"`
}

// Register function save user information into table users
func Register(username, email, password string) error {
	if !db.Select("id").Where("email = ?", email).First(&User{Email: email}).RecordNotFound() {
		return fmt.Errorf("email already registered")
	}
	if !db.Select("id").Where("username = ?", username).First(&User{Username: username}).RecordNotFound() {
		return fmt.Errorf("username already in use")
	}

	// TODO log error
	hash, _ := generatePassword(password)
	user := User{
		Email:    email,
		Username: username,
		Password: hash,
		Avatar:   gavatar(email, 80),
	}

	db.Create(&user)
	return nil
}

// Login function return user information
func Login(email, password string) (*User, error) {
	user := new(User)
	if db.Where("email = ?", email).First(user).RecordNotFound() {
		return nil, fmt.Errorf("email has not registerd")
	}

	if !verifyPassword(user.Password, password) {
		return nil, fmt.Errorf("password is incorrect")
	}

	return user, nil
}

func generatePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func gavatar(email string, size int) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(email)))
	return fmt.Sprintf("https://s.gravatar.com/avatar/%s?s=%d", hash, size)
}

// ComicSlice function get all comics subscribed by this user
func (u *User) ComicSlice() []*Comic {
	var comics []*Comic
	db.Model(u).Related(&comics, "Comics")
	return comics
}
