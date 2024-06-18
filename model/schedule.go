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

type Schedule struct {
	Base
	Name         string     `json:"name"`
	ScheduleType string     `json:"schedule_type" gorm:"type:ENUM('WEEKLY','DATERANGE','SINGLE_DATE')"`
	WeekDay      *string    `json:"week_day" gorm:"type:ENUM('SUNDAY','MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY')"`
	StartDate    *time.Time `json:"start_date" gorm:"type:DATE" sql:"TYPE:DATE"`
	EndDate      *time.Time `json:"end_date" gorm:"type:DATE" sql:"TYPE:DATE"`
	StartTime    *TimeOnly  `json:"start_time" gorm:"type:TIME"`
	EndTime      *TimeOnly  `json:"end_time" gorm:"type:TIME"`
	Employees    []Employee `json:"-" gorm:"many2many:schedule_employees;"`
	Sunday       bool       `json:"sunday"`
	Monday       bool       `json:"monday"`
	Tuesday      bool       `json:"tuesday"`
	Wednesday    bool       `json:"wednesday"`
	Thursday     bool       `json:"thursday"`
	Friday       bool       `json:"friday"`
	Saturday     bool       `json:"saturday"`
	EmployeeIDs  []string   `json:"employee_ids" gorm:"-"`
	CompanyID    string     `json:"company_id"`
	Company      Company    `gorm:"foreignKey:CompanyID"`
}

func (u *Schedule) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	return
}

func (m Schedule) MarshalJSON() ([]byte, error) {
	weekDay := ""
	startDate := ""
	endDate := ""
	startTime := ""
	endTime := ""
	if m.WeekDay != nil {
		weekDay = *m.WeekDay
	}
	if m.StartDate != nil {
		startDate = m.StartDate.Format("2006-01-02")
	}
	if m.EndDate != nil {
		endDate = m.EndDate.Format("2006-01-02")
	}
	if m.StartTime != nil {
		startTime = m.StartTime.Format("15:04")
	}
	if m.EndTime != nil {
		endTime = m.EndTime.Format("15:04")
	}

	employees := []resp.EmployeeResponse{}
	for _, v := range m.Employees {
		var picture *string
		var dateOfBirth, startedWork *time.Time
		if picture = &v.Picture.String; !v.Picture.Valid {
			picture = nil
		}
		if dateOfBirth = &v.DateOfBirth.Time; !v.DateOfBirth.Valid {
			dateOfBirth = nil
		}
		if startedWork = &v.StartedWork.Time; !v.StartedWork.Valid {
			startedWork = nil
		}
		var pictureUrl string
		if v.Picture.Valid {
			pictureUrl = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, v.Picture.String)
		}
		employees = append(employees, resp.EmployeeResponse{
			ID:                        v.ID,
			Email:                     v.Email,
			FirstName:                 v.FirstName,
			MiddleName:                v.MiddleName,
			LastName:                  v.LastName,
			Username:                  v.Username,
			Phone:                     v.Phone,
			JobTitle:                  v.JobTitle.Name,
			JobTitleID:                v.JobTitleID.String,
			Grade:                     v.Grade,
			Address:                   v.Address,
			Picture:                   picture,
			Cover:                     v.Cover,
			DateOfBirth:               dateOfBirth,
			EmployeeIdentityNumber:    v.EmployeeIdentityNumber,
			FullName:                  v.FullName,
			BasicSalary:               v.BasicSalary,
			PositionalAllowance:       v.PositionalAllowance,
			TransportAllowance:        v.TransportAllowance,
			MealAllowance:             v.MealAllowance,
			NonTaxableIncomeLevelCode: v.NonTaxableIncomeLevelCode,
			TaxPayerNumber:            v.TaxPayerNumber,
			Gender:                    v.Gender,
			OrganizationName:          v.Organization.Name,
			StartedWork:               startedWork,
			PictureUrl:                pictureUrl,
		})
	}

	return json.Marshal(resp.ScheduleResponse{
		ID:           m.ID,
		Name:         m.Name,
		ScheduleType: m.ScheduleType,
		WeekDay:      weekDay,
		StartDate:    startDate,
		EndDate:      endDate,
		StartTime:    startTime,
		EndTime:      endTime,
		Employees:    employees,
		Sunday:       m.Sunday,
		Monday:       m.Monday,
		Tuesday:      m.Tuesday,
		Wednesday:    m.Wednesday,
		Thursday:     m.Thursday,
		Friday:       m.Friday,
		Saturday:     m.Saturday,
	})
}
