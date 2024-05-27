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

func IncentiveSettingGetAllHandler(c *gin.Context) {
	var data []model.IncentiveSetting
	preloads := []string{"Shop", "ProductCategory"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }

	productCategoryId, ok := c.GetQuery("product_category_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"incentive_settings.product_category_id": productCategoryId,
		})

	}
	shopId, ok := c.GetQuery("shop_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"incentive_settings.shop_id": shopId,
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
		headers := []string{"No", "Toko", "Kategori", "T. Penjualan Min", "T. Penjualan Maks", "Komisi Min", "Komisi Min", "Batas Sakit", "Batas Izin", "Batas Alpa"}
		headerStyle := []int{xlsStyle.Bold, xlsStyle.Bold, xlsStyle.Bold,
			xlsStyle.TextRightBold, xlsStyle.TextRightBold, xlsStyle.TextRightBold, xlsStyle.TextRightBold,
			xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter,
		}
		headerWidth := []float64{7, 15, 15, 15, 15, 15, 15, 15, 15, 15}
		for j, v := range headers {
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), v)
			xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), headerStyle[j])
			xls.SetColWidth(sheet1Name, util.IntToLetters(int32(j)+1), util.IntToLetters(int32(j)+1), headerWidth[j])
		}
		row++
		for i, v := range data {
			cols := []interface{}{i + 1, v.Shop.Name, v.ProductCategory.Name, v.MinimumSalesTarget, v.MaximumSalesTarget, fmt.Sprintf("%v%s", v.MinimumSalesCommission*100, "%"), fmt.Sprintf("%v%s", v.MaximumSalesCommission*100, "%"), v.SickLeaveThreshold, v.OtherLeaveThreshold, v.AbsentThreshold}
			colStyle := []int{xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal,
				xlsStyle.TextRight, xlsStyle.TextRight, xlsStyle.TextRight, xlsStyle.TextRight,
				xlsStyle.Center, xlsStyle.Center, xlsStyle.Center,
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
		filename := fmt.Sprintf("Data-Setting-Insentif-%s.xlsx", time.Now().UTC().Format("02-01-2006"))
		c.Header("Content-Description", filename)
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b.Bytes())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List IncentiveSetting Retrived", dataRecords.Records, dataRecords)
}

func IncentiveSettingGetOneHandler(c *gin.Context) {
	var data model.IncentiveSetting

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveSetting Retrived", data, nil)
}

func IncentiveSettingCreateHandler(c *gin.Context) {
	var data model.IncentiveSetting

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
	util.ResponseSuccess(c, "Data IncentiveSetting Created", gin.H{"last_id": data.ID}, nil)
}

func IncentiveSettingUpdateHandler(c *gin.Context) {
	var input, data model.IncentiveSetting
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
	util.ResponseSuccess(c, "Data IncentiveSetting Updated", nil, nil)
}

func IncentiveSettingDeleteHandler(c *gin.Context) {
	var input, data model.IncentiveSetting
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveSetting Deleted", nil, nil)
}

func IncentiveImportHandler(c *gin.Context) {
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

		shop := model.Shop{}
		if err := database.DB.Find(&shop, "code = ?", row[2]).Error; err != nil {
			errString := fmt.Sprintf("Error at Line %s : %s", row[2], err.Error())
			errorRows = append(errorRows, errString)
			continue
		}

		cat := model.ProductCategory{}
		if err := database.DB.Find(&cat, "name = ?", row[3]).Error; err != nil {
			errString := fmt.Sprintf("Error at Line %s : %s", row[3], err.Error())
			errorRows = append(errorRows, errString)
			continue
		}
		minCommission := util.ExtractPercentage(row[6]) / 100
		maxCommission := util.ExtractPercentage(row[7]) / 100

		dataSetting := model.IncentiveSetting{
			ShopID:                 shop.ID,
			ProductCategoryID:      cat.ID,
			MinimumSalesTarget:     util.ParseThousandSeparatedNumber(row[4]),
			MaximumSalesTarget:     util.ParseThousandSeparatedNumber(row[5]),
			MinimumSalesCommission: minCommission,
			MaximumSalesCommission: maxCommission,
		}
		err = database.DB.Create(&dataSetting).Error
		if err != nil {
			errString := fmt.Sprintf("Error at create data %s : %s", row[2], err.Error())
			errorRows = append(errorRows, errString)
			continue
		}
		// util.LogJson(dataSetting)
		// for _, cell := range row {
		// 	fmt.Print(cell, "\t")
		// }
		// fmt.Println()
	}

	util.ResponseSuccess(c, "File imported", gin.H{
		"errors": errorRows,
	}, nil)
}
