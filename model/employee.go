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
	JobTitleID                *string
	JobTitle                  JobTitle `gorm:"foreignKey:JobTitleID"`
	Grade                     string
	Address                   string
	Picture                   *string
	Cover                     string
	StartedWork               *time.Time
	DateOfBirth               *time.Time
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
	OrganizationID            *string
}

func (u *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("id", uuid.New().String())
	names := []string{u.FirstName, u.MiddleName, u.LastName}
	tx.Statement.SetColumn("full_name", strings.Join(names, " "))

	return
}

func (m Employee) MarshalJSON() ([]byte, error) {

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
		Picture:                   m.Picture,
		Cover:                     m.Cover,
		DateOfBirth:               m.DateOfBirth,
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
		StartedWork:               m.StartedWork,
	})
}
