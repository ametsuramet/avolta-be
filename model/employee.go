package model

import (
	"avolta/config"
	"avolta/object/resp"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	Base
	Email                     string         `json:"email"`
	FirstName                 string         `json:"first_name"`
	MiddleName                string         `json:"middle_name"`
	LastName                  string         `json:"last_name"`
	Username                  string         `json:"username"`
	Phone                     string         `json:"phone"`
	JobTitleID                sql.NullString `json:"job_title_id"`
	JobTitle                  JobTitle       `gorm:"foreignKey:JobTitleID"`
	Grade                     string         `json:"grade"`
	Address                   string         `json:"address"`
	Picture                   sql.NullString `json:"picture"`
	Cover                     string         `json:"cover"`
	StartedWork               sql.NullTime   `json:"started_work"`
	DateOfBirth               sql.NullTime   `json:"date_of_birth"`
	EmployeeIdentityNumber    string         `json:"employee_identity_number"`
	FullName                  string         `json:"full_name"`
	ConnectedTo               sql.NullString `json:"connected_to"`
	Flag                      bool           `json:"flag"`
	BasicSalary               float64        `json:"basic_salary"`
	PositionalAllowance       float64        `json:"positional_allowance"`
	TransportAllowance        float64        `json:"transport_allowance"`
	MealAllowance             float64        `json:"meal_allowance"`
	NonTaxableIncomeLevelCode string         `json:"non_taxable_income_level_code"`
	PayRolls                  []PayRoll      `json:"pay_rolls"`
	TaxPayerNumber            string         `json:"tax_payer_number"`
	Gender                    string         `json:"gender"`
	Attendance                []Attendance   `json:"attendance" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Organization              Organization   `gorm:"foreignKey:OrganizationID"`
	OrganizationID            sql.NullString `json:"organization_id"`
	WorkingType               string         `gorm:"type:ENUM('FULL_TIME','PART_TIME','FREELANCE', 'FLEXIBLE','SHIFT','SEASONAL') DEFAULT 'FULL_TIME'"`
	Schedules                 []*Schedule    `json:"-" gorm:"many2many:schedule_employees;"`
}

func (u *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("id", uuid.New().String())
	}
	if u.FullName == "" {
		names := []string{u.FirstName, u.MiddleName, u.LastName}
		tx.Statement.SetColumn("full_name", strings.Join(names, " "))
	} else {
		splitName := strings.Split(u.FullName, " ")
		tx.Statement.SetColumn("first_name", splitName[0])
		if len(splitName) > 1 {
			tx.Statement.SetColumn("middle_name", splitName[1])
		}
		if len(splitName) > 2 {
			tx.Statement.SetColumn("last_name", splitName[2])
		}
	}
	return
}

func (m Employee) MarshalJSON() ([]byte, error) {
	var picture *string
	var dateOfBirth, startedWork *time.Time
	if picture = &m.Picture.String; !m.Picture.Valid {
		picture = nil
	}
	if dateOfBirth = &m.DateOfBirth.Time; !m.DateOfBirth.Valid {
		dateOfBirth = nil
	}
	if startedWork = &m.StartedWork.Time; !m.StartedWork.Valid {
		startedWork = nil
	}
	var pictureUrl string
	if m.Picture.Valid {
		pictureUrl = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, m.Picture.String)
	}
	schedules := []resp.ScheduleReponse{}
	for _, v := range m.Schedules {
		fmt.Println("Schedules", *v)
		if v != nil {
			weekDay := ""
			startDate := ""
			endDate := ""
			startTime := ""
			endTime := ""
			if v.WeekDay != nil {
				weekDay = *v.WeekDay
			}
			if v.StartDate != nil {
				startDate = v.StartDate.Format("2006-01-02")
			}
			if v.EndDate != nil {
				endDate = v.EndDate.Format("2006-01-02")
			}
			if v.StartTime != nil {
				startTime = v.StartTime.Format("15:04")
			}
			if v.EndTime != nil {
				endTime = v.EndTime.Format("15:04")
			}
			schedules = append(schedules, resp.ScheduleReponse{
				Name:         v.Name,
				ScheduleType: v.ScheduleType,
				WeekDay:      weekDay,
				StartDate:    startDate,
				EndDate:      endDate,
				StartTime:    startTime,
				EndTime:      endTime,
				Sunday:       v.Sunday,
				Monday:       v.Monday,
				Tuesday:      v.Tuesday,
				Wednesday:    v.Wednesday,
				Thursday:     v.Thursday,
				Friday:       v.Friday,
				Saturday:     v.Saturday,
			})
		}
	}

	return json.Marshal(resp.EmployeeReponse{
		ID:                        m.ID,
		Email:                     m.Email,
		FirstName:                 m.FirstName,
		MiddleName:                m.MiddleName,
		LastName:                  m.LastName,
		Username:                  m.Username,
		Phone:                     m.Phone,
		JobTitle:                  m.JobTitle.Name,
		JobTitleID:                m.JobTitleID.String,
		Grade:                     m.Grade,
		Address:                   m.Address,
		Picture:                   picture,
		Cover:                     m.Cover,
		DateOfBirth:               dateOfBirth,
		EmployeeIdentityNumber:    m.EmployeeIdentityNumber,
		FullName:                  m.FullName,
		BasicSalary:               m.BasicSalary,
		PositionalAllowance:       m.PositionalAllowance,
		TransportAllowance:        m.TransportAllowance,
		MealAllowance:             m.MealAllowance,
		NonTaxableIncomeLevelCode: m.NonTaxableIncomeLevelCode,
		TaxPayerNumber:            m.TaxPayerNumber,
		Gender:                    m.Gender,
		OrganizationName:          m.Organization.Name,
		StartedWork:               startedWork,
		PictureUrl:                pictureUrl,
		Schedules:                 schedules,
	})
}
