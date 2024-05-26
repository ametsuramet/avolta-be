package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/object/constants"
	"avolta/util"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ProductGetAllHandler(c *gin.Context) {
	var data []model.Product
	preloads := []string{"ProductCategory"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	search, ok := c.GetQuery("search")
	if ok {
		paginator.Search = append(paginator.Search, map[string]interface{}{
			"name":    search,
			"sku":     search,
			"barcode": search,
		})
	}
	productCategoryId, ok := c.GetQuery("product_category_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"product_category_id": productCategoryId,
		})

	}
	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	_, ok = c.GetQuery("download")
	if ok {
		sheet1Name := "Sheet1"
		row := 1
		xls := excelize.NewFile()
		xlsStyle := constants.NewExcelStyle(xls)
		xls.SetSheetName(xls.GetSheetName(0), sheet1Name)
		headers := []string{"No", "Nama Produk", "SKU", "Barcode", "Kategori", "Harga"}
		headerStyle := []int{xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold, xlsStyle.TextRightBold}
		headerWidth := []float64{7, 25, 15, 15, 15, 20, 30}
		for j, v := range headers {
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), v)
			xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), headerStyle[j])
			xls.SetColWidth(sheet1Name, util.IntToLetters(int32(j)+1), util.IntToLetters(int32(j)+1), headerWidth[j])
		}
		row++
		for i, v := range data {
			cols := []interface{}{i + 1, v.Name, v.SKU, v.Barcode, v.ProductCategory.Name, v.SellingPrice}
			colStyle := []int{xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.TextRight}
			for k, v := range cols {
				xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), v)
				xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), colStyle[k])
			}
			row++
		}

		var b *bytes.Buffer
		b, err := xls.WriteToBuffer()
		if err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
		filename := fmt.Sprintf("Data-Produk-%s.xlsx", time.Now().UTC().Format("02-01-2006"))
		c.Header("Content-Description", filename)
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b.Bytes())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List Product Retrived", dataRecords.Records, dataRecords)
}

func ProductGetOneHandler(c *gin.Context) {
	var data model.Product

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Product Retrived", data, nil)
}

func ProductCreateHandler(c *gin.Context) {
	var data model.Product

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// getUser, _ := c.Get("user")
	// user := getUser.(model.User)

	// data.AuthorID = user.ID

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Product Created", gin.H{"last_id": data.ID}, nil)
}

func ProductUpdateHandler(c *gin.Context) {
	var input, data model.Product
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Product Updated", nil, nil)
}

func ProductDeleteHandler(c *gin.Context) {
	var input, data model.Product
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Product Deleted", nil, nil)
}

func ProductImportHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()
	f, err := excelize.OpenReader(file)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	errorRows := []string{}
	rows, err := f.GetRows(f.WorkBook.Sheets.Sheet[0].Name)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	for i, row := range rows {
		if i < 1 {
			continue
		}
		cat := model.ProductCategory{}
		if err := database.DB.Find(&cat, "name = ?", row[4]).Error; err == nil {
			cat.Name = row[4]
			database.DB.Create(&cat)
		}
		price, err := strconv.Atoi(row[5])
		if err != nil {
			price = 0
		}
		if err := database.DB.Create(&model.Product{
			Name:              row[1],
			SKU:               row[2],
			Barcode:           row[3],
			ProductCategoryID: &cat.ID,
			SellingPrice:      float64(price),
		}).Error; err != nil {
			errString := fmt.Sprintf("Error at Line %s : %s", row[0], err.Error())
			errorRows = append(errorRows, errString)
		}
		for _, cell := range row {
			fmt.Print(cell, "\t")
		}
		fmt.Println()
	}

	util.ResponseSuccess(c, "File imported", gin.H{
		"errors": errorRows,
	}, nil)
}
