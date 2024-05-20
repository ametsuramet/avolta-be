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
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
)

func EmployeeImportHandler(c *gin.Context) {
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
		if i < 2 {
			continue
		}
		jobTitle := model.JobTitle{}
		database.DB.Find(&jobTitle, "name = ?", row[4])
		var dob, startedWork time.Time
		if dob, err = time.Parse("2006-01-02", row[6]); err != nil {
			dob = time.Time{}
		}
		if startedWork, err = time.Parse("2006-01-02", row[15]); err != nil {
			startedWork = time.Time{}
		}
		basicSalary, _ := strconv.Atoi(row[8])
		positionalAllowance, _ := strconv.Atoi(row[9])
		transportAllowance, _ := strconv.Atoi(row[10])
		mealAllowance, _ := strconv.Atoi(row[11])
		if err := database.DB.Create(&model.Employee{
			FullName:                  row[1],
			Email:                     row[2],
			Phone:                     row[3],
			JobTitleID:                model.NullStringConv(jobTitle.ID),
			Address:                   row[5],
			DateOfBirth:               model.NullTimeConv(dob),
			EmployeeIdentityNumber:    row[7],
			BasicSalary:               float64(basicSalary),
			PositionalAllowance:       float64(positionalAllowance),
			TransportAllowance:        float64(transportAllowance),
			MealAllowance:             float64(mealAllowance),
			NonTaxableIncomeLevelCode: row[12],
			TaxPayerNumber:            row[13],
			Gender:                    row[14],
			StartedWork:               model.NullTimeConv(startedWork),
		}).Error; err != nil {
			errString := fmt.Sprintf("Error at Line %s : %s", row[0], err.Error())
			errorRows = append(errorRows, errString)
		}
	}

	util.ResponseSuccess(c, "File imported", gin.H{
		"errors": errorRows,
	}, nil)
}
func EmployeeGetAllHandler(c *gin.Context) {
	var data []model.Employee
	preloads := []string{"JobTitle"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	paginator.Paginate(&data)
	search, ok := c.GetQuery("search")

	if ok {
		paginator.Search = append(paginator.Search, map[string]interface{}{
			"full_name":                search,
			"email":                    search,
			"employee_identity_number": search,
		})
	}
	jobTitleId, ok := c.GetQuery("job_title_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"job_title_id": jobTitleId,
		})

	}
	gender, ok := c.GetQuery("gender")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"gender": gender,
		})

	}
	ageStartDate, ok := c.GetQuery("age_start_date")
	if ok {
		paginator.WhereLessEqual = append(paginator.WhereLessEqual, map[string]interface{}{
			"date_of_birth": ageStartDate,
		})

	}
	ageEndDate, ok := c.GetQuery("age_end_date")
	if ok {
		paginator.WhereMoreEqual = append(paginator.WhereMoreEqual, map[string]interface{}{
			"date_of_birth": ageEndDate,
		})

	}
	startedWork, ok := c.GetQuery("started_work")
	if ok {
		paginator.WhereMoreEqual = append(paginator.WhereMoreEqual, map[string]interface{}{
			"started_work": startedWork,
		})

	}
	startedWorkEnd, ok := c.GetQuery("started_work_end")
	if ok {
		paginator.WhereLessEqual = append(paginator.WhereLessEqual, map[string]interface{}{
			"started_work": startedWorkEnd,
		})

	}
	_, ok = c.GetQuery("download")
	if ok {
		paginator.Preloads = append(paginator.Preloads, "JobTitle")

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
		headers := []string{"No", "Nama Karyawan", "NIK", "Jabatan", "Telp", "Email", "Alamat", "Jenis Kelamin", "Mulai Bekerja"}
		headerStyle := []int{xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter, xlsStyle.BoldCenter}
		headerWidth := []float64{7, 25, 15, 15, 15, 20, 30, 15, 15}
		for j, v := range headers {
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), v)
			xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(j)+1), row), headerStyle[j])
			xls.SetColWidth(sheet1Name, util.IntToLetters(int32(j)+1), util.IntToLetters(int32(j)+1), headerWidth[j])
		}
		row++
		orderNumber := 1
		for _, v := range data {
			var genderStr string
			var startedWorkStr string
			if genderStr = "Laki-laki"; v.Gender == "f" {
				genderStr = "Perempuan"
			}
			if startedWorkStr = ""; v.StartedWork.Valid {
				startedWorkStr = v.StartedWork.Time.Format("02-01-2006")
			}
			cols := []interface{}{orderNumber, v.FullName, v.EmployeeIdentityNumber, v.JobTitle.Name, v.Phone, v.Email, v.Address, genderStr, startedWorkStr}
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
		filename := fmt.Sprintf("Data-Karyawan-%s.xlsx", time.Now().UTC().Format("02-01-2006"))
		c.Header("Content-Description", filename)
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b.Bytes())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List Employee Retrived", dataRecords.Records, dataRecords)
}

func EmployeeGetOneHandler(c *gin.Context) {
	var data model.Employee

	id := c.Params.ByName("id")

	if err := database.DB.Preload("Schedules").Preload("JobTitle").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Employee Retrived", data, nil)
}

func EmployeeCreateHandler(c *gin.Context) {
	var data model.Employee

	if err := c.ShouldBindWith(&data, binding.JSON); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// getUser, _ := c.Get("user")
	// user := getUser.(model.User)

	// data.AuthorID = user.ID

	util.LogJson(data)

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Employee Created", gin.H{"last_id": data.ID}, nil)
}

func EmployeeUpdateHandler(c *gin.Context) {
	var input, data model.Employee
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
	util.ResponseSuccess(c, "Data Employee Updated", nil, nil)
}

func EmployeeDeleteHandler(c *gin.Context) {
	var input, data model.Employee
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Employee Deleted", nil, nil)
}
