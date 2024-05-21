package model

import (
	"avolta/config"
	"avolta/database"
	"avolta/object/resp"
	"avolta/util"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	Base
	ClockIn                time.Time            `json:"clock_in"`
	ClockOut               *time.Time           `json:"clock_out"`
	ClockInNotes           string               `json:"clock_in_notes"`
	ClockOutNotes          string               `json:"clock_out_notes"`
	ClockInPicture         string               `json:"clock_in_picture"`
	ClockOutPicture        string               `json:"clock_out_picture"`
	ClockInLat             float64              `json:"clock_in_lat" gorm:"type:DECIMAL(10,8)"`
	ClockInLng             float64              `json:"clock_in_lng" gorm:"type:DECIMAL(11,8)"`
	ClockOutLat            float64              `json:"clock_out_lat" gorm:"type:DECIMAL(10,8)"`
	ClockOutLng            float64              `json:"clock_out_lng" gorm:"type:DECIMAL(11,8)"`
	EmployeeID             *string              `json:"employee_id"`
	Employee               Employee             `gorm:"foreignKey:EmployeeID"`
	BreakStart             *TimeOnly            `json:"break_start" gorm:"type:TIME"`
	BreakEnd               *TimeOnly            `json:"break_end" gorm:"type:TIME"`
	Overtime               *TimeOnly            `json:"overtime" gorm:"type:TIME"`
	LateIn                 *TimeOnly            `json:"late_in" gorm:"type:TIME"`
	WorkingDuration        *TimeOnly            `json:"working_duration" gorm:"type:TIME"`
	AttendanceBulkImportID *string              `json:"attendance_bulk_import_id"`
	AttendanceBulkImport   AttendanceBulkImport `gorm:"foreignKey:AttendanceBulkImportID"`
	AttendanceImportItemID *string              `json:"attendance_import_item_id"`
	AttendanceImportItem   AttendanceImportItem `gorm:"foreignKey:AttendanceImportItemID"`
}

func (u *Attendance) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (u *Attendance) AfterUpdate(tx *gorm.DB) (err error) {
	if u.ClockOut != nil {
		emp := Employee{}
		database.DB.Select("daily_working_hours").Find(&emp, "id = ?", u.EmployeeID)

		scannedOverTime, err := time.Parse("15:04", util.FormatDuration(u.ClockOut.Sub(u.ClockIn)-(time.Duration(emp.DailyWorkingHours)*time.Hour)))
		if err == nil {
			overTime := &TimeOnly{scannedOverTime}
			tx.Model(&Attendance{}).Where("id = ?", u.ID).Update("overtime", overTime)
		}
		scannedDuration, err := time.Parse("15:04", util.FormatDuration(u.ClockOut.Sub(u.ClockIn)))
		if err == nil {
			duration := &TimeOnly{scannedDuration}
			tx.Model(&Attendance{}).Where("id = ?", u.ID).Update("working_duration", duration)
		}
	}
	return
}

func (m Attendance) MarshalJSON() ([]byte, error) {
	var employeePicture string
	if m.Employee.Picture.Valid {
		employeePicture = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, m.Employee.Picture.String)
	}
	overTime := ""
	if m.Overtime != nil {
		overTime = m.Overtime.Format("15:04")
	}
	duration := ""
	if m.WorkingDuration != nil {
		duration = m.WorkingDuration.Format("15:04")
	}

	return json.Marshal(resp.AttendanceReponse{
		ID:                     m.ID,
		ClockIn:                m.ClockIn,
		ClockOut:               m.ClockOut,
		ClockInNotes:           m.ClockInNotes,
		ClockOutNotes:          m.ClockOutNotes,
		ClockInPicture:         m.ClockInPicture,
		ClockOutPicture:        m.ClockOutPicture,
		ClockInLat:             m.ClockInLat,
		ClockInLng:             m.ClockInLng,
		ClockOutLat:            m.ClockOutLat,
		ClockOutLng:            m.ClockOutLng,
		EmployeeName:           m.Employee.FullName,
		EmployeeID:             m.EmployeeID,
		EmployeeJobTitle:       m.Employee.JobTitle.Name,
		EmployeePicture:        &employeePicture,
		EmployeeIdentityNumber: m.Employee.EmployeeIdentityNumber,
		Overtime:               overTime,
		WorkingDuration:        duration,
	})
}
