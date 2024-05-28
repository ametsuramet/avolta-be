package handler

import (
	"avolta/config"
	"avolta/database"
	"avolta/model"
	"avolta/object/constants"
	"avolta/util"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TigorLazuardi/tanggal"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func AttendanceImportRejectHandler(c *gin.Context) {
	var data model.AttendanceBulkImport
	var input = struct {
		Notes string `json:"notes"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	data.Notes = input.Notes
	data.Status = "REJECTED"
	if err := database.DB.Updates(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data Attendance Bulk Import Rejected", data, nil)
}

func AttendanceImportApproveHandler(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var data model.AttendanceBulkImport
	var input = struct {
		Notes string `json:"notes"`
	}{}

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Params.ByName("id")

	if err := database.DB.Preload("Data", func(db *gorm.DB) *gorm.DB {
		db = db.Order("sequence_number asc")
		return db
	}).Preload("Data.Items", func(db *gorm.DB) *gorm.DB {
		db = db.Order("sequence_number asc")
		return db
	}).Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, dataAttendances := range data.Data {
			var employee model.Employee
			if err := database.DB.Find(&employee, "id = ?", dataAttendances.SystemEmployeeID).Error; err != nil {
				return err
			}
			for _, item := range dataAttendances.Items {
				if item.DutyOn == "" {
					continue
				}
				if item.DutyOn == "" && item.DutyOff == "" {
					continue
				}
				clockIn, err := time.ParseInLocation("02-01-2006 15:04", fmt.Sprintf("%s %s", item.Date, item.DutyOn), loc)
				if err != nil {
					return err
				}
				var clockOut *time.Time
				if item.DutyOff != "" {
					cOut, err := time.ParseInLocation("02-01-2006 15:04", fmt.Sprintf("%s %s", item.Date, item.DutyOff), loc)
					if err != nil {
						return err
					}
					clockOut = &cOut
				}

				var overTime *model.TimeOnly
				if item.Overtime != "" {
					scannedTime, err := time.Parse("15:04", item.Overtime)
					if err == nil {
						overTime = &model.TimeOnly{scannedTime}
					}
				}
				var lateIn *model.TimeOnly
				if item.LateIn != "" {
					scannedTime, err := time.Parse("15:04", item.LateIn)
					if err == nil {
						lateIn = &model.TimeOnly{scannedTime}
					}
				}

				if err := database.DB.Create(&model.Attendance{
					ClockIn:                clockIn,
					ClockOut:               clockOut,
					Overtime:               overTime,
					ClockInNotes:           item.Notes,
					EmployeeID:             &employee.ID,
					AttendanceBulkImportID: &id,
					AttendanceImportItemID: &item.ID,
					LateIn:                 lateIn,
				}).Error; err != nil {
					return err
				}

			}

		}

		data.Status = "APPROVED"
		data.Notes = input.Notes
		database.DB.Updates(&data)

		return nil
	})
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Attendance Bulk Import Approved", nil, nil)
}

func AttendanceImportDetailHandler(c *gin.Context) {
	var data model.AttendanceBulkImport

	id := c.Params.ByName("id")

	if err := database.DB.Preload("User").Preload("Data", func(db *gorm.DB) *gorm.DB {
		db = db.Order("sequence_number asc")
		return db
	}).Preload("Data.Items", func(db *gorm.DB) *gorm.DB {
		db = db.Order("sequence_number asc")
		return db
	}).Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Attendance Bulk Import Retrived", data, nil)

}
func AttendanceImportHandler(c *gin.Context) {
	file, headers, err := c.Request.FormFile("file")
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

	fmt.Println(errorRows)
	seqNo := 1
	dataMappings := []model.AttendanceImport{}
	dataMapping := model.AttendanceImport{}
	seqItemNo := 1
	for _, row := range rows {
		if len(row) == 0 {
			continue

		}
		// for _, cell := range row {
		// 	fmt.Print(cell, "\t")
		// 	fmt.Println()
		// }

		if row[0] == "Fingerprint ID" {
			dataMapping.SequenceNumber = seqNo
			dataMapping.FingerprintID = row[2]
			dataMapping.Items = []model.AttendanceImportItem{}

		}
		if row[0] == "Employee Code" {
			dataMapping.EmployeeCode = row[2]
			employee := model.Employee{}
			database.DB.Find(&employee, "employee_code = ?", row[2])
			dataMapping.SystemEmployeeName = employee.FullName
			dataMapping.SystemEmployeeID = employee.ID
		}
		if row[0] == "Employee Name" {
			dataMapping.EmployeeName = row[2]
		}
		if util.Contains(config.DaysOfWeekAbbr, row[0]) {
			dataMapping.Items = append(dataMapping.Items, model.AttendanceImportItem{
				SequenceNumber: seqItemNo,
				Day:            row[0],
				Date:           row[1],
				WorkingHour:    row[2],
				Activity:       row[3],
				DutyOn:         row[4],
				DutyOff:        row[5],
				LateIn:         row[6],
				EarlyDeparture: row[7],
				EffectiveHour:  row[8],
				Overtime:       row[9],
				Notes:          row[10],
			})
			seqItemNo++
		}

		if len(row) > 3 {
			if row[3] == "Total" {
				dataMappings = append(dataMappings, dataMapping)
				seqNo++
				seqItemNo = 1
			}
		}

	}
	getUser, _ := c.Get("user")
	user := getUser.(model.User)

	data := model.AttendanceBulkImport{
		FileName:       headers.Filename,
		ImportedBy:     user.ID,
		DateImportedAt: time.Now(),
		Data:           dataMappings,
		Status:         "DRAFT",
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		database.DB.Create(&data)
		database.DB.Save(&data)
		return nil
	})
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Import Succeed", data.ID, nil)
}

func SummaryGetAllHandler(c *gin.Context) {
	var data model.Employee

	employeeId := c.Params.ByName("employeeId")

	if err := database.DB.Find(&data, "id = ?", employeeId).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	startDate, ok := c.GetQuery("start_date")
	if !ok {
		util.ResponseFail(c, http.StatusBadRequest, "Rentang tanggal belum di set")
		return
	}
	endDate, ok := c.GetQuery("end_date")
	if !ok {
		util.ResponseFail(c, http.StatusBadRequest, "Rentang tanggal belum di set")
		return
	}

	summaries := []struct {
		ID       string  `sql:"id"`
		ClockIn  string  `sql:"clock_in"`
		ClockOut string  `sql:"clock_out"`
		Diff     float64 `sql:"diff"`
		FullTime float64 `sql:"full_time"`
	}{}

	if err := database.DB.Raw(`select id, clock_in , clock_out , SUM(TIMESTAMPDIFF(MINUTE ,clock_in, clock_out))  as diff,  
	CASE WHEN SUM(TIMESTAMPDIFF(MINUTE ,clock_in, clock_out)) / 60 >= ? THEN 1 ELSE 0 END full_time
	from attendances a 
	WHERE DATE(clock_in) BETWEEN DATE(?) AND DATE(?)
	AND clock_in is not null
	and employee_id = ?
	group by id, DATE(clock_in), DATE(clock_out) `, data.DailyWorkingHours, startDate, endDate, employeeId).Scan(&summaries).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	totalDays := 0
	totalFullDays := 0
	totalMinutes := 0
	for _, v := range summaries {
		totalDays++
		if v.FullTime == 1 {
			totalFullDays++
		}
		totalMinutes += int(v.Diff)
	}

	util.ResponsePaginatorSuccess(c, "Data Summary Attandance Retrived", gin.H{
		"employee_name":       data.FullName,
		"total_working_days":  data.TotalWorkingDays,
		"total_working_hours": data.TotalWorkingHours,
		"daily_working_hours": data.DailyWorkingHours,
		"total_hours":         totalMinutes,
		"total_days":          totalDays,
		"total_full_days":     totalFullDays,
	}, nil)

}
func AttendanceGetAllHandler(c *gin.Context) {
	_, ok := c.GetQuery("download")
	if ok {
		b, err := generateReport(c)
		if err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
		startDate, _ := c.GetQuery("start_date")
		endDate, _ := c.GetQuery("end_date")

		startDateFormat, _ := time.Parse(time.RFC3339, startDate)
		endDateFormat, _ := time.Parse(time.RFC3339, endDate)

		filename := fmt.Sprintf("Laporan-Absensi-%s-%s.xlsx", startDateFormat.Format("02-01-2006"), endDateFormat.Format("02-01-2006"))
		c.Header("Content-Description", filename)
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", b)
		return
	}
	var data []model.Attendance

	preloads := []string{"Employee", "Employee.JobTitle"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	search, ok := c.GetQuery("search")

	paginator.Joins = append(paginator.Joins, map[string]interface{}{
		"LEFT JOIN employees ON employees.id = attendances.employee_id": nil,
	})

	if ok {
		paginator.Search = append(paginator.Search, map[string]interface{}{
			"employees.full_name":                search,
			"employees.employee_identity_number": search,
		})
	}
	jobTitleId, ok := c.GetQuery("job_title_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"employees.job_title_id": jobTitleId,
		})

	}
	gender, ok := c.GetQuery("gender")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"employees.gender": gender,
		})

	}
	startDate, ok := c.GetQuery("start_date")
	if ok {
		paginator.WhereMoreEqual = append(paginator.WhereMoreEqual, map[string]interface{}{
			"attendances.clock_in": startDate,
		})

	}
	endDate, ok := c.GetQuery("end_date")
	if ok {
		paginator.WhereLessEqual = append(paginator.WhereLessEqual, map[string]interface{}{
			"attendances.clock_in": endDate,
		})

	}
	employeeIds, ok := c.GetQuery("employee_ids")
	if ok {
		paginator.WhereIn = append(paginator.WhereIn,
			map[string][]string{
				"attendances.employee_id": strings.Split(employeeIds, ","),
			},
		)

	}
	employeeId, ok := c.GetQuery("employee_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"attendances.employee_id": employeeId,
		})

	}

	paginator.OrderBy = []string{"clock_in asc"}
	orderBy, ok := c.GetQuery("order_by")
	if ok {
		paginator.OrderBy = []string{orderBy}

	}

	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponsePaginatorSuccess(c, "Data List Attandance Retrived", dataRecords.Records, dataRecords)
}

func generateReport(c *gin.Context) ([]byte, error) {
	timezone, _ := c.Get("timezone")
	loc, _ := time.LoadLocation(timezone.(string))
	time.Local = loc
	employeeIds, ok := c.GetQuery("employee_ids")
	if !ok {
		return []byte{}, errors.New("no employee id")
	}
	startDate, ok := c.GetQuery("start_date")
	if !ok {
		return []byte{}, errors.New("no start date")
	}
	endDate, ok := c.GetQuery("end_date")
	if !ok {
		return []byte{}, errors.New("no end date")
	}
	xls := excelize.NewFile()
	xlsStyle := constants.NewExcelStyle(xls)

	format := []tanggal.Format{
		tanggal.Hari,
		tanggal.NamaBulan,
		tanggal.Tahun,
	}

	for index, v := range strings.Split(employeeIds, ",") {
		// var data []model.Attendance
		// database.DB.Find(&data, "employee_id = ? and clock_in >= ? and clock_out <= ?", v, startDate, endDate)
		var employee model.Employee
		database.DB.Find(&employee, "id = ?", v)
		row := 1

		startDateFormat, _ := time.ParseInLocation(time.RFC3339, startDate, loc)
		endDateFormat, _ := time.ParseInLocation(time.RFC3339, endDate, loc)

		dateList := util.GetDates(startDateFormat, endDateFormat)

		fmt.Println("DATELIST", dateList)
		sheet1Name := employee.FullName
		if index == 0 {
			xls.SetSheetName(xls.GetSheetName(index), sheet1Name)
		} else {
			xls.NewSheet(sheet1Name)
		}
		xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", "A", row), "Perhitungan Jam Kerja")
		xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "A", row), fmt.Sprintf("%s%v", "B", row), xlsStyle.BgBlueTextWhite)
		xls.MergeCell(sheet1Name, fmt.Sprintf("%s%v", "A", row), fmt.Sprintf("%s%v", "B", row))
		xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", "D", row), "Ringkasan")
		xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "D", row), fmt.Sprintf("%s%v", "E", row), xlsStyle.BgBlueTextWhite)
		xls.MergeCell(sheet1Name, fmt.Sprintf("%s%v", "D", row), fmt.Sprintf("%s%v", "E", row))
		xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", "G", row), "Keterlambatan")
		xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "G", row), fmt.Sprintf("%s%v", "H", row), xlsStyle.BgBlueTextWhite)
		xls.MergeCell(sheet1Name, fmt.Sprintf("%s%v", "G", row), fmt.Sprintf("%s%v", "H", row))
		row++
		xls.SetColWidth(sheet1Name, "A", "A", 20)
		xls.SetColWidth(sheet1Name, "B", "B", 20)
		xls.SetColWidth(sheet1Name, "C", "C", 20)
		xls.SetColWidth(sheet1Name, "D", "D", 20)
		xls.SetColWidth(sheet1Name, "E", "E", 20)
		xls.SetColWidth(sheet1Name, "F", "F", 20)
		xls.SetColWidth(sheet1Name, "G", "G", 17.5)
		xls.SetColWidth(sheet1Name, "H", "H", 17.5)
		xls.SetColWidth(sheet1Name, "I", "I", 17.5)
		xls.SetColWidth(sheet1Name, "J", "J", 17.5)
		startDateLocale, _ := tanggal.Papar(startDateFormat, "", tanggal.WIB)
		endDateLocale, _ := tanggal.Papar(endDateFormat, "", tanggal.WIB)
		monthLabel := startDateLocale.Format(" ", []tanggal.Format{tanggal.NamaBulan, tanggal.Tahun})
		if startDateFormat.Month() != endDateFormat.Month() {
			monthLabel = startDateLocale.Format(" ", []tanggal.Format{tanggal.NamaBulan, tanggal.Tahun}) + " ~ " + endDateLocale.Format(" ", []tanggal.Format{tanggal.NamaBulan, tanggal.Tahun})
		}

		groupRow := [][]string{
			{"Nama Karyawan", employee.FullName, "", "Hadir", ""},
			{"NIP/NIK", employee.EmployeeIdentityNumber, "", "Ijin", ""},
			{"Bulan", monthLabel, "", "Sakit", ""},
			{"Periode", fmt.Sprintf("%s ~ %s", startDateLocale.Format(" ", format), endDateLocale.Format(" ", format)), "", "Dinas", ""},
			{"Jumlah Hari Kerja", "", "", "Cuti", ""},
			{"Jumlah Jam Kerja Wajib", "", "", "Remote", ""},
			{"Lembur ", "", "", "Alfa", ""},
		}
		for _, rowData := range groupRow {
			for k, v := range rowData {
				xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), v)
				xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "A", row), fmt.Sprintf("%s%v", "A", row), xlsStyle.Bold)
				xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "D", row), fmt.Sprintf("%s%v", "D", row), xlsStyle.Bold)
			}
			row++
		}
		row++
		xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", "A", row), "Daily Activity")
		xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "A", row), fmt.Sprintf("%s%v", "A", row), xlsStyle.GreenPastel)
		row++
		headers := []string{"No", "Tgl", "Hari", "Jam Masuk", "Jam Keluar", "Jumlah Jam Istirahat", "Jumlah Jam Kerja", "Kehadiran", "Keterlambatan", "Aktifitas"}
		for k, v := range headers {
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), v)
			xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(k)+1), row), xlsStyle.GreenPastel)
		}
		row++
		numb := 1
		firstData := row
		for _, date := range dateList {
			attendances := []model.Attendance{}
			database.DB.Order("clock_in asc").Find(&attendances, "employee_id = ? and DATE(clock_in) = ?", v, date.Local().Format("2006-01-02"))
			tgl, _ := tanggal.Papar(date, "", tanggal.WIB)
			dayFormat := tgl.Format(" ", []tanggal.Format{tanggal.NamaHari})
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(1), row), numb)
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(2), row), tgl.Format(" ", format))
			xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(3), row), dayFormat)
			xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(1), row), fmt.Sprintf("%s%v", util.IntToLetters(3), row), xlsStyle.Center)
			if dayFormat == "Minggu" {
				xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(1), row), fmt.Sprintf("%s%v", util.IntToLetters(3), row), xlsStyle.CenterRed)
			}
			for indexRowData, rowData := range attendances {

				var _, hourIn, hourOut, hourBreak, hourWork, att, late, activity interface{}
				if rowData.ID != "" {
					hourIn, hourOut, hourBreak, hourWork, att, late, activity = rowData.ClockIn.Format("15:04"), rowData.ClockOut.Format("15:04"), "", rowData.ClockOut.Sub(rowData.ClockIn).Hours(), "", "", ""
				}
				if indexRowData > 0 {
					xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(1), row), numb)
					xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(1), row), fmt.Sprintf("%s%v", util.IntToLetters(1), row), xlsStyle.Center)
				} else {
					if len(attendances) > 0 {
						xls.MergeCell(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(2), row), fmt.Sprintf("%s%v", util.IntToLetters(2), row+len(attendances)-1))
						xls.MergeCell(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(3), row), fmt.Sprintf("%s%v", util.IntToLetters(3), row+len(attendances)-1))
						xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(2), row), fmt.Sprintf("%s%v", util.IntToLetters(3), row+len(attendances)-1), xlsStyle.CenterCenter)
					}
				}

				cells := []interface{}{hourIn, hourOut, hourBreak, hourWork, att, late, activity}
				for m, v := range cells {
					xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(m)+4), row), v)
					if m == 3 {
						xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(m)+4), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(m)+4), row), xlsStyle.NumberCenter)
					} else {
						xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", util.IntToLetters(int32(m)+4), row), fmt.Sprintf("%s%v", util.IntToLetters(int32(m)+4), row), xlsStyle.Center)

					}
				}
				row++
				numb++
			}
			if len(attendances) == 0 {
				row++
				numb++
			}

		}
		xls.SetCellValue(sheet1Name, fmt.Sprintf("%s%v", "C", row), "Total")
		xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "C", row), fmt.Sprintf("%s%v", "C", row), xlsStyle.BoldCenter)

		xls.SetCellFormula(sheet1Name, fmt.Sprintf("%s%v", "G", row), fmt.Sprintf("SUM(%s:%s)", fmt.Sprintf("%s%v", "G", firstData), fmt.Sprintf("%s%v", "G", row-1)))
		xls.SetCellStyle(sheet1Name, fmt.Sprintf("%s%v", "G", row), fmt.Sprintf("%s%v", "G", row), xlsStyle.NumberCenter)

		row++
	}

	var b *bytes.Buffer
	b, err := xls.WriteToBuffer()
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return []byte{}, errors.New("error generate file")
	}
	return b.Bytes(), nil
}

func AttendanceGetOneHandler(c *gin.Context) {
	var data model.Attendance

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Attendance Retrived", data, nil)
}

func AttendanceCreateHandler(c *gin.Context) {
	var data model.Attendance

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
	util.ResponseSuccess(c, "Data Attendance Created", gin.H{"last_id": data.ID}, nil)
}

func AttendanceUpdateHandler(c *gin.Context) {
	var input, data model.Attendance
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Preload("Employee").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	// if input.ClockOut != nil {

	// 	scannedOverTime, err := time.Parse("15:04", util.FormatDuration(input.ClockOut.Sub(input.ClockIn)-(time.Duration(data.Employee.DailyWorkingHours)*time.Hour)))
	// 	if err == nil {
	// 		input.Overtime = &model.TimeOnly{scannedOverTime}
	// 	}
	// 	scannedDuration, err := time.Parse("15:04", util.FormatDuration(input.ClockOut.Sub(input.ClockIn)))
	// 	if err == nil {
	// 		input.WorkingDuration = &model.TimeOnly{scannedDuration}
	// 	}
	// }
	if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Attendance Updated", nil, nil)
}

func AttendanceDeleteHandler(c *gin.Context) {
	var input, data model.Attendance
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Attendance Deleted", nil, nil)
}
