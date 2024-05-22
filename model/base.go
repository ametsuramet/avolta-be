package model

import (
	"avolta/config"
	"avolta/database"
	"avolta/util"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Base struct {
	ID        string         `gorm:"type:char(36);primary_key" json:"id"`
	CreatedAt time.Time      `json:"-" `
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `sql:"index" json:"-"`
}
type BaseNoTimeStamp struct {
	ID string `gorm:"type:char(36);primary_key" json:"id"`
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

func NullStringConv(str string) sql.NullString {
	if len(str) > 0 {
		return sql.NullString{String: str, Valid: true}
	}
	return sql.NullString{String: str, Valid: false}
}
func NullTimeConv(check time.Time) sql.NullTime {
	zero := time.Time{}
	if check == zero {
		return sql.NullTime{Time: zero, Valid: false}
	}
	return sql.NullTime{Time: check, Valid: true}
}

func ExtractNumber(data Setting, number string) string {
	re2 := regexp.MustCompile(`{(.*?)}`)
	pattern := "0.{1}"
	if data.PayRollAutoNumberCharacterLength > 0 {
		pattern = fmt.Sprintf("(\\d{%d})", data.PayRollAutoNumberCharacterLength)

	}
	re := regexp.MustCompile(pattern)
	getNumber := re.FindAllString(number, -1)
	if len(getNumber) == 1 {
		num, _ := strconv.Atoi(getNumber[0])
		if num > 0 {
			return GenerateInvoiceBillNumber(data, getNumber[0])
		}
	}
	values := []any{}
	for _, v := range re2.FindAllStringSubmatch(data.PayRollAutoFormat, -1) {
		if len(v) > 0 {

			if v[1] == config.STATIC_CHARACTER {
				values = append(values, data.PayRollStaticCharacter)
			} else if v[1] == config.AUTO_NUMERIC {
				values = append(values, "(\\d+)")
			} else if v[1] == config.MONTH_ROMAN {
				intMonth, _ := strconv.Atoi(time.Now().Format("1"))
				values = append(values, util.IntegerToRoman(intMonth))
			} else if v[1] == config.MONTH_MM {
				values = append(values, time.Now().Format("01"))
			} else if v[1] == config.MONTH_MMM {
				values = append(values, time.Now().Format("Jan"))
			} else if v[1] == config.MONTH_MMMM {
				values = append(values, time.Now().Format("January"))
			} else if v[1] == config.YEAR_YY {
				values = append(values, time.Now().Format("06"))
			} else if v[1] == config.YEAR_YYYY {
				values = append(values, time.Now().Format("2006"))
			} else {
				values = append(values, v[0])
			}
		}
	}
	pattern2 := strings.ReplaceAll(fmt.Sprintf(re2.ReplaceAllString(data.PayRollAutoFormat, "%s"), values...), "/", "\\/")
	re3 := regexp.MustCompile(pattern2)

	for _, v := range re3.FindAllStringSubmatch(number, -1) {
		if len(v) > 0 {
			return GenerateInvoiceBillNumber(data, v[1])
		}
	}

	return GenerateInvoiceBillNumber(data, "00")
}
func GenerateInvoiceBillNumber(data Setting, before string) string {

	re := regexp.MustCompile(`{(.*?)}`)
	values := []any{}
	for _, v := range re.FindAllStringSubmatch(data.PayRollAutoFormat, -1) {
		if len(v) > 0 {
			if v[1] == config.STATIC_CHARACTER {
				values = append(values, data.PayRollStaticCharacter)
			} else if v[1] == config.AUTO_NUMERIC {

				numberBefore, err := strconv.Atoi(before)
				if err != nil {
					values = append(values, before)
				} else {
					if data.PayRollAutoNumberCharacterLength == 0 {
						values = append(values, fmt.Sprintf("%d", numberBefore+1))
					} else {
						length := strconv.Itoa(data.PayRollAutoNumberCharacterLength)

						values = append(values, fmt.Sprintf("%0"+length+"d", numberBefore+1))
					}

				}
			} else if v[1] == config.MONTH_ROMAN {
				intMonth, _ := strconv.Atoi(time.Now().Format("1"))
				values = append(values, util.IntegerToRoman(intMonth))
			} else if v[1] == config.MONTH_MM {
				values = append(values, time.Now().Format("01"))
			} else if v[1] == config.MONTH_MMM {
				values = append(values, time.Now().Format("Jan"))
			} else if v[1] == config.MONTH_MMMM {
				values = append(values, time.Now().Format("January"))
			} else if v[1] == config.YEAR_YY {
				values = append(values, time.Now().Format("06"))
			} else if v[1] == config.YEAR_YYYY {
				values = append(values, time.Now().Format("2006"))
			} else {
				values = append(values, v[0])
			}
		}
	}

	return fmt.Sprintf(re.ReplaceAllString(data.PayRollAutoFormat, "%s"), values...)
}
