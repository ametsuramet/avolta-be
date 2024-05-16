package model

import (
	"avolta/database"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Base struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

// func (u *Base) BeforeCreate(tx *gorm.DB) (err error) {
// 	tx.Statement.SetColumn("id", uuid.New().String())

// 	return
// }

// paginate retrieves paginated records for a given model
func Paginate(c *gin.Context, model interface{}, preloads []string, args ...interface{}) (int64, error) {
	count := int64(0)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return count, errors.New("invalid page number")

	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		return count, errors.New("invalid limit value")
	}

	order := c.DefaultQuery("order", "created_at desc")

	// Calculate the offset based on the page and limit
	offset := (page - 1) * limit

	database.DB.Model(model).Select("id").Find(model, args...).Count(&count)

	// Retrieve records with pagination
	result := database.DB.Model(model).Order(order)
	if len(preloads) > 0 {
		for _, v := range preloads {
			result = result.Preload(v)
		}
	}
	result = result.Offset(offset).Limit(limit).Find(model, args...)
	if result.Error != nil {
		return count, fmt.Errorf("failed to retrieve records for %T", model)
	}

	return count, nil
}
