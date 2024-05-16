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
	Email                     string         `json:"email"`
	FirstName                 string         `json:"first_name"`
	MiddleName                string         `json:"middle_name"`
	LastName                  string         `json:"last_name"`
	Username                  string         `json:"username"`
	Phone                     string         `json:"phone"`
	Position                  string         `json:"position"`
	Grade                     string         `json:"grade"`
	Address                   string         `json:"address"`
	Picture                   *string        `json:"picture"`
	Cover                     string         `json:"cover"`
	StartedWork               *time.Time     `json:"started_work"`
	DateOfBirth               *time.Time     `json:"date_of_birth"`
	EmployeeIdentityNumber    string         `json:"employee_identity_number"`
	FullName                  string         `json:"full_name"`
	ConnectedTo               sql.NullString `json:"connected_to" `
	Flag                      bool           `json:"flag"`
	BasicSalary               float64        `json:"basic_salary"`
	PositionalAllowance       float64        `json:"positional_allowance"`
	TransportAllowance        float64        `json:"transport_allowance"`
	MealAllowance             float64        `json:"meal_allowance"`
	NonTaxableIncomeLevelCode string         `json:"non_taxable_income_level_code"`
	PayRolls                  []PayRoll      `json:"pay_rolls"`
	TaxPayerNumber            string         `json:"tax_payer_number"`
	Gender                    string         `json:"gender" gorm:"type:enum('f', 'm')"`
	Attendance                []Attendance   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Organization              Organization   `gorm:"foreignKey:OrganizationID"`
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
		Position:                  m.Position,
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
