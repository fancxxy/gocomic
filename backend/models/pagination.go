package models

import (
	"math"

	"github.com/jinzhu/gorm"
)

// paginate function
func paginate(limit, page int) *gorm.DB {
	limit = int(math.Max(1, float64(limit)))
	page = int(math.Max(1, float64(page)))
	return db.Offset(limit * (page - 1)).Limit(limit)
}
