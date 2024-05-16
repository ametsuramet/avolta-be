package model

import (
	"avolta/object/resp"
	"encoding/json"
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
	Employee        Employee `gorm:"foreignKey:EmployeeID"`
}

func (u *Attendance) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	return
}

func (m Attendance) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.AttendanceReponse{
		ClockIn:          m.ClockIn,
		ClockOut:         m.ClockOut,
		ClockInNotes:     m.ClockInNotes,
		ClockOutNotes:    m.ClockOutNotes,
		ClockInPicture:   m.ClockInPicture,
		ClockOutPicture:  m.ClockOutPicture,
		ClockInLat:       m.ClockInLat,
		ClockInLng:       m.ClockInLng,
		ClockOutLat:      m.ClockOutLat,
		ClockOutLng:      m.ClockOutLng,
		EmployeeName:     m.Employee.FullName,
		EmployeeID:       m.EmployeeID,
		EmployeePosition: m.Employee.Position,
		EmployeePicture:  m.Employee.Picture,
	})
}
