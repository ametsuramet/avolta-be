package model

import (
	"avolta/config"
	"avolta/object/resp"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	Base
	ClockIn         time.Time
	ClockOut        *time.Time
	ClockInNotes    string
	ClockOutNotes   string
	ClockInPicture  string
	ClockOutPicture string
	ClockInLat      float64 `gorm:"type:DECIMAL(10,8)"`
	ClockInLng      float64 `gorm:"type:DECIMAL(11,8)"`
	ClockOutLat     float64 `gorm:"type:DECIMAL(10,8)"`
	ClockOutLng     float64 `gorm:"type:DECIMAL(11,8)"`
	EmployeeID      *string
	Employee        Employee  `gorm:"foreignKey:EmployeeID"`
	BreakStart      *TimeOnly `gorm:"type:TIME"`
	BreakEnd        *TimeOnly `gorm:"type:TIME"`
	Overtime        *TimeOnly `gorm:"type:TIME"`
}

func (u *Attendance) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Attendance) MarshalJSON() ([]byte, error) {
	var employeePicture string
	if m.Employee.Picture.Valid {
		employeePicture = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, m.Employee.Picture.String)
	}

	return json.Marshal(resp.AttendanceReponse{
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
	})
}
