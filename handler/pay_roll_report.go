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

func PayRollReportGetAllHandler(c *gin.Context) {
	var data []model.PayRollReport
	preloads := []string{"User"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }
	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List PayRollReport Retrived", dataRecords.Records, dataRecords)
}

func PayRollReportGetOneHandler(c *gin.Context) {
	var data model.PayRollReport

	id := c.Params.ByName("id")

	if err := database.DB.Preload("Items").Preload("User").Preload("Items.PayRoll").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRollReport Retrived", data, nil)
}

func PayRollReportCreateHandler(c *gin.Context) {
	var data model.PayRollReport

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	getUser, _ := c.Get("user")
	user := getUser.(model.User)

	data.UserID = user.ID

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRollReport Created", gin.H{"last_id": data.ID}, nil)
}

func PayRollReportUpdateHandler(c *gin.Context) {
	var input, data model.PayRollReport
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
	util.ResponseSuccess(c, "Data PayRollReport Updated", nil, nil)
}
func PayRollReporDownloadPayRollBankHandler(c *gin.Context) {
	var data model.PayRollReport
	id := c.Params.ByName("id")

	if err := database.DB.Preload("Items").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	sheet1Name := "Sheet1"
	row := 1
	xls := excelize.NewFile()
	xlsStyle := constants.NewExcelStyle(xls)
	xls.SetSheetName(xls.GetSheetName(0), sheet1Name)
	headers := []string{"No", "Nama Karyawan", "Telp", "Email", "Bank", "Kode Bank", "No. Rekening"}
	headerStyle := []int{xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter}
	headerWidth := []float64{7, 25, 15, 15, 15, 20, 30, 15, 15}
	for j, v := range headers {
		xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), v)
		xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), headerStyle[j])
		xls.SetColWidth(sheet1Name, util.IntToLetters(int32(j)+1), util.IntToLetters(int32(j)+1), headerWidth[j])
	}
	row++
	orderNumber := 1

	for _, v := range data.Items {
		cols := []interface{}{orderNumber, v.EmployeeName, v.EmployeePhone, v.EmployeeEmail, v.BankName, v.BankCode, v.BankAccountNumber}
		colStyle := []int{xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal, xlsStyle.Normal}
		for k, v := range cols {
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), v)
			xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), colStyle[k])
		}
		row++
		orderNumber++
	}

	var b *bytes.Buffer
	b, err := xls.WriteToBuffer()
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	filename := fmt.Sprintf("PayRoll-Bank-%s-%s.xlsx", data.ReportNumber, time.Now().UTC().Format("02-01-2006"))
	c.Header("Content-Description", filename)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b.Bytes())
}
func PayRollReportAddItemHandler(c *gin.Context) {
	var data model.PayRollReport
	id := c.Params.ByName("id")
	payrollID := c.Params.ByName("payrollID")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := data.AddItem(payrollID); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data PayRollReport Add Item", nil, nil)
}
func PayRollReportDeleteItemHandler(c *gin.Context) {
	var data model.PayRollReport
	id := c.Params.ByName("id")
	payrollID := c.Params.ByName("payrollID")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := data.DeleteItem(payrollID); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data PayRollReport Delete Item", nil, nil)
}
func PayRollReportDeleteHandler(c *gin.Context) {
	var input, data model.PayRollReport
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRollReport Deleted", nil, nil)
}
