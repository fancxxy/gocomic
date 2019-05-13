package routers

import "github.com/fancxxy/gocomic/backend/models"

// Register contains register information
type Register struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login contains login information
type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Subscribe contains login information
type Subscribe struct {
	Website string `json:"website" binding:"required"`
	Comic   string `json:"comic" binding:"required"`
}

// Unsubscribe contains login information
type Unsubscribe struct {
	ID int `json:"id" binding:"required"`
}

// Response struct
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type comicData struct {
	Time  models.Time   `json:"subscribed"`
	Comic *models.Comic `json:"comic"`
}

type comicsData struct {
	Comics []*models.Comic `json:"comics"`
}

type chapterData struct {
	Chapter *models.Chapter `json:"chapter"`
}
