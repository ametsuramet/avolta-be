package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/object/constants"
	"avolta/util"
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func SaleGetAllHandler(c *gin.Context) {
	var data []model.Sale
	preloads := []string{"Product", "Product.ProductCategory", "Employee", "Shop"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	paginator.Joins = append(paginator.Joins, map[string]interface{}{
		"JOIN products ON products.id = sales.product_id": nil,
	})
	paginator.Joins = append(paginator.Joins, map[string]interface{}{
		"JOIN employees ON employees.id = sales.employee_id": nil,
	})

	search, ok := c.GetQuery("search")
	if ok {
		paginator.Search = append(paginator.Search, map[string]interface{}{
			"sales.code":       search,
			"products.name":    search,
			"products.sku":     search,
			"products.barcode": search,
		})
	}

	startDate, ok := c.GetQuery("start_date")
	if ok {
		paginator.WhereMoreEqual = append(paginator.WhereMoreEqual, map[string]interface{}{
			"sales.date": startDate,
		})

	}
	endDate, ok := c.GetQuery("end_date")
	if ok {
		paginator.WhereLessEqual = append(paginator.WhereLessEqual, map[string]interface{}{
			"sales.date": endDate,
		})

	}
	productCategoryId, ok := c.GetQuery("product_category_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"products.product_category_id": productCategoryId,
		})

	}
	shopId, ok := c.GetQuery("shop_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"sales.shop_id": shopId,
		})

	}
	employeeId, ok := c.GetQuery("employee_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"sales.employee_id": employeeId,
		})

	}
	incentiveId, ok := c.GetQuery("incentive_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"sales.incentive_id": incentiveId,
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
		// ac := accounting.Accounting{Symbol: "", Precision: 4}
		// no	date	name	sku	barcode	qty	price	sub_total	discount	discount_amount	total	sales	nik	shop	shop_code
		headers := []string{"No", "Tgl", "Nama Produk", "SKU", "Kategori", "Qty", "Sub Total", "Diskon", "Diskon Amount", "Total", "Salesman", "Nik", "Toko", "Kode Toko"}
		headerStyle := []int{xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold,
			xlsStyle.TextRightBold, xlsStyle.TextRightBold, xlsStyle.TextRightBold, xlsStyle.TextRightBold, xlsStyle.TextRightBold,
			xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold,
		}
		headerWidth := []float64{7, 25, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}
		for j, v := range headers {
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), v)
			xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), headerStyle[j])
			xls.SetColWidth(sheet1Name, util.IntToLetters(int32(j)+1), util.IntToLetters(int32(j)+1), headerWidth[j])
		}
		row++
		for i, v := range data {
			cols := []interface{}{i + 1, v.Date.Format("02-01-2006"), v.Product.Name, v.Product.SKU, v.Product.ProductCategory.Name,
				v.Qty, v.SubTotal, v.Discount * 100, v.DiscountAmount, v.Total, v.Employee.FullName, v.Employee.EmployeeIdentityNumber, v.Shop.Name, v.Shop.Code,
			}
			colStyle := []int{xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal,
				xlsStyle.TextRight, xlsStyle.TextRight, xlsStyle.TextRight, xlsStyle.TextRight, xlsStyle.TextRight,
				xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal,
			}
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
		filename := fmt.Sprintf("Data-Penjualan-%s.xlsx", time.Now().UTC().Format("02-01-2006"))
		c.Header("Content-Description", filename)
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b.Bytes())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List Sale Retrived", dataRecords.Records, dataRecords)
}

func SaleGetOneHandler(c *gin.Context) {
	var data model.Sale

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Sale Retrived", data, nil)
}

func SaleCreateHandler(c *gin.Context) {
	var data model.Sale

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
	util.ResponseSuccess(c, "Data Sale Created", gin.H{"last_id": data.ID}, nil)
}

func SaleUpdateHandler(c *gin.Context) {
	var input, data model.Sale
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
	util.ResponseSuccess(c, "Data Sale Updated", nil, nil)
}

func SaleDeleteHandler(c *gin.Context) {
	var input, data model.Sale
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Sale Deleted", nil, nil)
}

func SaleImportHandler(c *gin.Context) {
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
		product := model.Product{}
		if err := database.DB.Find(&product, "sku = ?", row[3]).Error; err != nil {
			errString := fmt.Sprintf("Error at Line %s : %s", row[3], err.Error())
			errorRows = append(errorRows, errString)
			continue
		}
		employee := model.Employee{}
		if err := database.DB.Find(&employee, "employee_identity_number = ?", row[12]).Error; err != nil {
			errString := fmt.Sprintf("Error at Line %s : %s", row[12], err.Error())
			errorRows = append(errorRows, errString)
			continue
		}
		shop := model.Shop{}
		if err := database.DB.Find(&shop, "code = ?", row[14]); err != nil {

			shop.Name = row[13]
			shop.Code = row[14]
			database.DB.Create(&shop)
			// continue
		}
		date, err := time.Parse("02-01-2006", row[1])
		if err != nil {
			errString := fmt.Sprintf("Error at Line %s : %s", row[1], err.Error())
			errorRows = append(errorRows, errString)
			continue
		}

		qty := util.ParseThousandSeparatedNumber(row[5])

		price := util.ParseThousandSeparatedNumber(row[6])

		subTotal := util.ParseThousandSeparatedNumber(row[7])

		discount := util.ExtractPercentage(row[8]) / 100
		discountAmount := util.ParseThousandSeparatedNumber(row[9])

		total := util.ParseThousandSeparatedNumber(row[10])

		data := model.Sale{
			Date:           date,
			ProductID:      product.ID,
			EmployeeID:     employee.ID,
			Discount:       discount,
			Qty:            qty,
			Price:          price,
			SubTotal:       subTotal,
			DiscountAmount: discountAmount,
			Total:          total,
			ShopID:         shop.ID,
		}
		err = database.DB.Create(&data).Error
		if err != nil {
			errString := fmt.Sprintf("Error at create data %s : %s", row[2], err.Error())
			errorRows = append(errorRows, errString)
			continue
		}

		util.LogJson(data)
		if i < 1 {
			continue
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
