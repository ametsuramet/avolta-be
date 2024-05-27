package util

import (
	"avolta/database"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Paginator struct {
	DB             *gorm.DB
	Ctx            *gin.Context
	OrderBy        []string
	Page           string
	PerPage        string
	Search         []map[string]interface{}
	Where          []map[string]interface{}
	WhereOr        []map[string]interface{}
	WhereQuery     []map[string][]interface{}
	Joins          []map[string]interface{}
	WhereNotNull   []string
	WhereNull      []string
	WhereLess      []map[string]interface{}
	WhereMore      []map[string]interface{}
	WhereLessEqual []map[string]interface{}
	WhereMoreEqual []map[string]interface{}
	WhereIn        []map[string][]string
	WhereNotIn     []map[string][]string
	Table          string
	Group          string
	Select         []string
	Preloads       []string
}

type DataPaginator struct {
	TotalRecords int64       `json:"total_records"`
	Records      interface{} `json:"records"`
	CurrentPage  int64       `json:"current_page"`
	TotalPages   int64       `json:"total_pages"`
	Next         int64       `json:"next"`
	Prev         int64       `json:"prev"`
}

// func (p *Paginator) WhereQuery(query interface{}, args ...interface{}) *gorm.DB {
// 	fmt.Println("Execute WhereQuery")
// 	db := p.DB
// 	db = db.Where(query, args...)
// 	p.DB = db

// 	return db
// }

func (p *Paginator) Paginate(model interface{}) (*DataPaginator, error) {
	done := make(chan bool, 1)
	if p.Ctx == nil {
		return nil, errors.New("no context loaded")
	}

	count := int64(0)
	page, err := strconv.Atoi(p.Ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return nil, errors.New("invalid page number")

	}

	limit, err := strconv.Atoi(p.Ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 0 {
		return nil, errors.New("invalid limit value")
	}

	p.OrderBy = append(p.OrderBy, p.Ctx.DefaultQuery("order", "created_at desc"))
	offset := (page - 1) * limit

	db := database.DB
	p.DB = db

	for _, o := range p.OrderBy {
		db = db.Order(o)
	}

	for _, where := range p.Where {
		for key, value := range where {
			db = db.Where(key+"=?", value)
		}
	}

	for _, join := range p.Joins {
		for key, value := range join {
			if value == nil {
				db = db.Joins(key)
			} else {
				db = db.Joins(key, value)
			}
		}
	}

	for _, where := range p.WhereIn {
		for key, value := range where {
			db = db.Where(key+" IN (?)", value)
		}
	}

	for _, where := range p.WhereNotIn {
		for key, value := range where {
			db = db.Where(key+" NOT IN (?)", value)
		}
	}

	for _, where := range p.WhereLess {
		for key, value := range where {
			db = db.Where(key+" < ?", value)
		}
	}

	for _, where := range p.WhereMore {
		for key, value := range where {
			db = db.Where(key+" > ?", value)
		}
	}
	for _, where := range p.WhereMoreEqual {
		for key, value := range where {
			db = db.Where(key+" >= ?", value)
		}
	}
	for _, where := range p.WhereLessEqual {
		for key, value := range where {
			db = db.Where(key+" <= ?", value)
		}
	}
	for _, value := range p.WhereNotNull {
		db = db.Where(value + " is not null")
	}
	for _, value := range p.WhereNull {
		db = db.Where(value + " is null")
	}
	for _, preload := range p.Preloads {
		db = db.Preload(preload)
	}
	for _, query := range p.WhereQuery {
		for key, value := range query {
			db = db.Where(key, value...)
		}
		// db = db.Preload(preload)
	}

	var searchQuery []string
	for _, search := range p.Search {
		var j = 0
		for key, value := range search {
			// fmt.Println(j, key, value)
			if j == 0 {
				searchQuery = append(searchQuery, key+" like \"%"+value.(string)+"%\"")
				// db = db.Where(key+" like ?", "%"+value.(string)+"%")
			} else {
				searchQuery = append(searchQuery, "OR "+key+" like \"%"+value.(string)+"%\"")
				// db = db.Or(key+" like ?", "%"+value.(string)+"%")
			}
			j++
		}
	}

	if len(p.Search) > 0 {
		var searchQueryRaw = "( " + strings.Join(searchQuery, " ") + " )"
		// fmt.Println(searchQueryRaw)
		db = db.Where(searchQueryRaw)
	}

	countRecords(db, model, done, &count)

	if len(p.Select) > 1 {
		db = db.Select(p.Select)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	db = db.Find(model)
	err = db.Error
	if err != nil {
		return nil, err
	}

	<-done
	totalPage := int64(0)
	if limit != 0 {
		totalPage = int64(math.Ceil(float64(count) / float64(limit)))
	}
	prev := int64(1)
	next := totalPage
	if page-1 != 0 {
		prev = int64(page) - 1
	}

	return &DataPaginator{
		TotalRecords: count,
		CurrentPage:  int64(page),
		Records:      model,
		TotalPages:   totalPage,
		Next:         next,
		Prev:         prev,
	}, nil
}

func countRecords(db *gorm.DB, countDataSource interface{}, done chan bool, count *int64) {
	db.Model(countDataSource).Count(count)
	done <- true
}

func NewPaginator(c *gin.Context) Paginator {
	return Paginator{
		Ctx: c,
	}
}
