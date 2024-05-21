package util

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ResponsePaginatorSuccess(c *gin.Context, msg string, data interface{}, pagination *DataPaginator) {
	resp := gin.H{
		"message": msg,
	}
	if data != nil {
		resp["data"] = data
	}

	if pagination != nil {
		resp["pagination"] = gin.H{
			"total_records": pagination.TotalRecords,
			"current_page":  pagination.CurrentPage,
			"total_pages":   pagination.TotalPages,
			"next":          pagination.Next,
			"prev":          pagination.Prev,
		}
	}

	c.JSON(http.StatusOK, resp)
}
func ResponseSuccess(c *gin.Context, msg string, data interface{}, totalRecords *int64) {
	resp := gin.H{
		"message": msg,
	}
	if data != nil {
		resp["data"] = data
	}
	if totalRecords != nil {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		totalPage := math.Ceil(float64(*totalRecords) / float64(limit))

		prev := 1
		next := totalPage
		if page-1 != 0 {
			prev = page - 1
		}
		resp["pagination"] = gin.H{
			"page":          page,
			"next":          next,
			"prev":          prev,
			"total_page":    totalPage,
			"total_records": *totalRecords,
		}
	}

	c.JSON(http.StatusOK, resp)
}

func ResponseFail(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg, "message": msg})
}
