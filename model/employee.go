package model

import (
	"avolta/object/resp"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	Base
	Email                     string
	FirstName                 string
	MiddleName                string
	LastName                  string
	Username                  string
	Phone                     string
	JobTitleID                sql.NullString
	JobTitle                  JobTitle `gorm:"foreignKey:JobTitleID"`
	Grade                     string
	Address                   string
	Picture                   sql.NullString
	Cover                     string
	StartedWork               sql.NullTime
	DateOfBirth               sql.NullTime
	EmployeeIdentityNumber    string
	FullName                  string
	ConnectedTo               sql.NullString
	Flag                      bool
	BasicSalary               float64
	PositionalAllowance       float64
	TransportAllowance        float64
	MealAllowance             float64
	NonTaxableIncomeLevelCode string
	PayRolls                  []PayRoll
	TaxPayerNumber            string
	Gender                    string
	Attendance                []Attendance `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Organization              Organization `gorm:"foreignKey:OrganizationID"`
	OrganizationID            sql.NullString
}

func (u *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
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

	return json.Marshal(resp.EmployeeReponse{
		ID:                        m.ID,
		Email:                     m.Email,
		FirstName:                 m.FirstName,
		MiddleName:                m.MiddleName,
		LastName:                  m.LastName,
		Username:                  m.Username,
		Phone:                     m.Phone,
		JobTitle:                  m.JobTitle.Name,
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
	})
}
