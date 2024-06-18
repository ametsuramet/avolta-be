package model

import (
	"avolta/config"
	"avolta/database"
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
	Email                     string          `json:"email"`
	FirstName                 string          `json:"first_name"`
	MiddleName                string          `json:"middle_name"`
	LastName                  string          `json:"last_name"`
	Username                  string          `json:"username"`
	Phone                     string          `json:"phone"`
	JobTitleID                sql.NullString  `json:"job_title_id"`
	JobTitle                  JobTitle        `gorm:"foreignKey:JobTitleID"`
	Grade                     string          `json:"grade"`
	UserID                    *string         `json:"user_id"`
	Address                   string          `json:"address"`
	Picture                   sql.NullString  `json:"picture"`
	Cover                     string          `json:"cover"`
	StartedWork               sql.NullTime    `json:"started_work"`
	DateOfBirth               sql.NullTime    `json:"date_of_birth"`
	EmployeeIdentityNumber    string          `json:"employee_identity_number"`
	EmployeeCode              string          `json:"employee_code"`
	FullName                  string          `json:"full_name"`
	ConnectedTo               sql.NullString  `json:"connected_to"`
	Flag                      bool            `json:"flag"`
	BasicSalary               float64         `json:"basic_salary"`
	PositionalAllowance       float64         `json:"positional_allowance"`
	TransportAllowance        float64         `json:"transport_allowance"`
	MealAllowance             float64         `json:"meal_allowance"`
	NonTaxableIncomeLevelCode string          `json:"non_taxable_income_level_code"`
	PayRolls                  []PayRoll       `json:"pay_rolls" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reimbursements            []Reimbursement `json:"reimbursements" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TaxPayerNumber            string          `json:"tax_payer_number"`
	Gender                    string          `json:"gender"`
	Attendance                []Attendance    `json:"attendance" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Organization              Organization    `gorm:"foreignKey:OrganizationID"`
	OrganizationID            sql.NullString  `json:"organization_id"`
	WorkingType               string          `json:"working_type" gorm:"type:ENUM('FULL_TIME','PART_TIME','FREELANCE', 'FLEXIBLE','SHIFT','SEASONAL') DEFAULT 'FULL_TIME'"`
	Schedules                 []*Schedule     `json:"-" gorm:"many2many:schedule_employees;"`
	TotalWorkingDays          int32           `json:"total_working_days"`
	TotalWorkingHours         float64         `json:"total_working_hours"`
	DailyWorkingHours         float64         `json:"daily_working_hours"`
	WorkSafetyRisks           string          `gorm:"type:ENUM('very_low','low','middle', 'high','very_high') DEFAULT 'very_low'" json:"work_safety_risks"`
	Sales                     []Sale          `json:"sales" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BankAccountNumber         string          `json:"bank_account_number"`
	BankID                    *string         `json:"bank_id"`
	Bank                      Bank            `gorm:"foreignKey:BankID"`
	CompanyID                 string          `json:"company_id" gorm:"not null"`
	Company                   Company         `gorm:"foreignKey:CompanyID"`
}

func (u *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	fmt.Println("CONTEXT", ctx)
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
func (u *Employee) BeforeUpdate(tx *gorm.DB) (err error) {
	// fmt.Println("UPDATE EMPLOYEE")
	if u.FullName == "" {
		// fmt.Println("UPDATE EMPLOYEE", 1)
		names := []string{u.FirstName, u.MiddleName, u.LastName}
		tx.Statement.SetColumn("full_name", strings.Join(names, " "))
	} else {
		// fmt.Println("UPDATE EMPLOYEE", 2)
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
	schedules := []resp.ScheduleResponse{}
	for _, v := range m.Schedules {
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
			schedules = append(schedules, resp.ScheduleResponse{
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
	username := ""
	userID := ""
	if m.UserID != nil {
		user := User{}
		database.DB.Select("full_name").Find(&user, "id = ?", m.UserID)
		username = user.FullName
		userID = *m.UserID
	}

	bankName := ""
	if m.BankID != nil {
		bankName = m.Bank.Name
	}

	return json.Marshal(resp.EmployeeResponse{
		ID:                        m.ID,
		Email:                     m.Email,
		FirstName:                 m.FirstName,
		MiddleName:                m.MiddleName,
		LastName:                  m.LastName,
		Username:                  username,
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
		EmployeeCode:              m.EmployeeCode,
		UserID:                    userID,
		TotalWorkingDays:          m.TotalWorkingDays,
		TotalWorkingHours:         m.TotalWorkingHours,
		DailyWorkingHours:         m.DailyWorkingHours,
		BankAccountNumber:         m.BankAccountNumber,
		BankID:                    m.BankID,
		BankName:                  bankName,
		WorkingType:               m.WorkingType,
	})
}

type NonTaxableIncomeLevel struct {
	Code                         string  `json:"code"`
	Description                  string  `json:"description"`
	Amount                       float64 `json:"amount"`
	EffectiveRateAverageCategory string  `json:"effective_rate_average_category"`
}

func NewNonTaxableIncomeLevel() []NonTaxableIncomeLevel {
	return []NonTaxableIncomeLevel{
		{"-", "Non Pajak", 0, ""},
		{"TK/0", "Tidak Kawin Tanpa Tanggungan", 54000000, "A"},
		{"TK/1", "Tidak Kawin 1 Orang Tanggungan", 58500000, "A"},
		{"TK/2", "Tidak Kawin 2 Orang Tanggungan", 63000000, "B"},
		{"TK/3", "Tidak Kawin 3 Orang Tanggungan", 67500000, "B"},
		{"K/0", "Kawin Tanpa Tanggungan", 58500000, "A"},
		{"K/1", "Kawin 1 Orang Tanggungan", 63000000, "B"},
		{"K/2", "Kawin 2 Orang Tanggungan", 67500000, "B"},
		{"K/3", "Kawin 3 Orang Tanggungan", 72000000, "C"},
		{"K/1/0", "Kawin Penghasilan Istri Digabung Dengan Suami Tanpa Tanggungan", 112500000, "A"},
		{"K/1/1", "Kawin Penghasilan Istri Digabung Dengan Suami 1 Orang Tanggungan", 117000000, "A"},
		{"K/1/2", "Kawin Penghasilan Istri Digabung Dengan Suami 2 Orang Tanggungan", 121500000, "B"},
		{"K/1/3", "Kawin Penghasilan Istri Digabung Dengan Suami 3 Orang Tanggungan", 126000000, "C"},
	}
}

func (m Employee) GetNonTaxableIncomeLevelAmount() float64 {
	for _, v := range NewNonTaxableIncomeLevel() {
		if v.Code == m.NonTaxableIncomeLevelCode {
			return v.Amount
		}
	}
	return 0
}
func (m Employee) GetNonTaxableIncomeLevelCategory() string {
	for _, v := range NewNonTaxableIncomeLevel() {
		if v.Code == m.NonTaxableIncomeLevelCode {
			return v.EffectiveRateAverageCategory
		}
	}
	return ""
}

func (m Employee) EffectiveRateAverageCategoryA(taxable float64) float64 {
	if taxable < 5400000 {
		return float64(0) / 100
	} else if taxable < 5650000 {
		return float64(0.25) / 100
	} else if taxable < 5950000 {
		return float64(0.5) / 100
	} else if taxable < 6300000 {
		return float64(0.75) / 100
	} else if taxable < 6750000 {
		return float64(1) / 100
	} else if taxable < 7500000 {
		return float64(1.25) / 100
	} else if taxable < 8550000 {
		return float64(1.5) / 100
	} else if taxable < 9650000 {
		return float64(1.75) / 100
	} else if taxable < 10050000 {
		return float64(2) / 100
	} else if taxable < 10350000 {
		return float64(2.25) / 100
	} else if taxable < 10700000 {
		return float64(2.5) / 100
	} else if taxable < 11050000 {
		return float64(3) / 100
	} else if taxable < 11600000 {
		return float64(3.5) / 100
	} else if taxable < 12500000 {
		return float64(4) / 100
	} else if taxable < 13750000 {
		return float64(5) / 100
	} else if taxable < 15100000 {
		return float64(6) / 100
	} else if taxable < 16950000 {
		return float64(7) / 100
	} else if taxable < 19750000 {
		return float64(8) / 100
	} else if taxable < 24150000 {
		return float64(9) / 100
	} else if taxable < 26450000 {
		return float64(10) / 100
	} else if taxable < 28000000 {
		return float64(11) / 100
	} else if taxable < 30050000 {
		return float64(12) / 100
	} else if taxable < 32400000 {
		return float64(13) / 100
	} else if taxable < 35400000 {
		return float64(14) / 100
	} else if taxable < 39100000 {
		return float64(15) / 100
	} else if taxable < 43850000 {
		return float64(16) / 100
	} else if taxable < 47800000 {
		return float64(17) / 100
	} else if taxable < 51400000 {
		return float64(18) / 100
	} else if taxable < 56300000 {
		return float64(19) / 100
	} else if taxable < 62200000 {
		return float64(20) / 100
	} else if taxable < 68600000 {
		return float64(21) / 100
	} else if taxable < 77500000 {
		return float64(22) / 100
	} else if taxable < 89000000 {
		return float64(23) / 100
	} else if taxable < 103000000 {
		return float64(24) / 100
	} else if taxable < 125000000 {
		return float64(25) / 100
	} else if taxable < 157000000 {
		return float64(26) / 100
	} else if taxable < 206000000 {
		return float64(27) / 100
	} else if taxable < 337000000 {
		return float64(28) / 100
	} else if taxable < 454000000 {
		return float64(29) / 100
	} else if taxable < 550000000 {
		return float64(30) / 100
	} else if taxable < 695000000 {
		return float64(31) / 100
	} else if taxable < 910000000 {
		return float64(32) / 100
	} else if taxable < 1400000000 {
		return float64(33) / 100
	} else {
		return float64(34) / 100
	}
}

func (m Employee) EffectiveRateAverageCategoryB(taxable float64) float64 {
	if taxable < 6200000 {
		return float64(0) / 100
	} else if taxable < 6500000 {
		return float64(0.25) / 100
	} else if taxable < 6850000 {
		return float64(0.5) / 100
	} else if taxable < 7300000 {
		return float64(0.75) / 100
	} else if taxable < 9200000 {
		return float64(1) / 100
	} else if taxable < 10750000 {
		return float64(1.5) / 100
	} else if taxable < 11250000 {
		return float64(2) / 100
	} else if taxable < 11600000 {
		return float64(2.5) / 100
	} else if taxable < 12600000 {
		return float64(3) / 100
	} else if taxable < 13600000 {
		return float64(4) / 100
	} else if taxable < 14950000 {
		return float64(5) / 100
	} else if taxable < 16400000 {
		return float64(6) / 100
	} else if taxable < 18450000 {
		return float64(7) / 100
	} else if taxable < 21850000 {
		return float64(8) / 100
	} else if taxable < 26000000 {
		return float64(9) / 100
	} else if taxable < 27700000 {
		return float64(10) / 100
	} else if taxable < 29350000 {
		return float64(11) / 100
	} else if taxable < 31450000 {
		return float64(12) / 100
	} else if taxable < 33950000 {
		return float64(13) / 100
	} else if taxable < 37100000 {
		return float64(14) / 100
	} else if taxable < 41100000 {
		return float64(15) / 100
	} else if taxable < 45800000 {
		return float64(16) / 100
	} else if taxable < 49500000 {
		return float64(17) / 100
	} else if taxable < 53800000 {
		return float64(18) / 100
	} else if taxable < 58500000 {
		return float64(19) / 100
	} else if taxable < 64000000 {
		return float64(20) / 100
	} else if taxable < 71000000 {
		return float64(21) / 100
	} else if taxable < 80000000 {
		return float64(22) / 100
	} else if taxable < 93000000 {
		return float64(23) / 100
	} else if taxable < 109000000 {
		return float64(24) / 100
	} else if taxable < 129000000 {
		return float64(25) / 100
	} else if taxable < 163000000 {
		return float64(26) / 100
	} else if taxable < 211000000 {
		return float64(27) / 100
	} else if taxable < 374000000 {
		return float64(28) / 100
	} else if taxable < 459000000 {
		return float64(29) / 100
	} else if taxable < 555000000 {
		return float64(30) / 100
	} else if taxable < 704000000 {
		return float64(31) / 100
	} else if taxable < 957000000 {
		return float64(32) / 100
	} else if taxable < 1405000000 {
		return float64(33) / 100
	} else {
		return float64(34) / 100
	}
}

func (m Employee) EffectiveRateAverageCategoryC(taxable float64) float64 {
	if taxable < 6600000 {
		return float64(0) / 100
	} else if taxable < 6950000 {
		return float64(0.25) / 100
	} else if taxable < 7350000 {
		return float64(0.5) / 100
	} else if taxable < 7800000 {
		return float64(0.75) / 100
	} else if taxable < 8850000 {
		return float64(1) / 100
	} else if taxable < 9800000 {
		return float64(1.25) / 100
	} else if taxable < 10950000 {
		return float64(1.5) / 100
	} else if taxable < 11200000 {
		return float64(1.75) / 100
	} else if taxable < 12050000 {
		return float64(2) / 100
	} else if taxable < 12950000 {
		return float64(3) / 100
	} else if taxable < 14150000 {
		return float64(4) / 100
	} else if taxable < 15550000 {
		return float64(5) / 100
	} else if taxable < 17050000 {
		return float64(6) / 100
	} else if taxable < 19500000 {
		return float64(7) / 100
	} else if taxable < 22700000 {
		return float64(8) / 100
	} else if taxable < 26600000 {
		return float64(9) / 100
	} else if taxable < 28100000 {
		return float64(10) / 100
	} else if taxable < 30100000 {
		return float64(11) / 100
	} else if taxable < 32600000 {
		return float64(12) / 100
	} else if taxable < 35400000 {
		return float64(13) / 100
	} else if taxable < 38900000 {
		return float64(14) / 100
	} else if taxable < 43000000 {
		return float64(15) / 100
	} else if taxable < 47400000 {
		return float64(16) / 100
	} else if taxable < 51200000 {
		return float64(17) / 100
	} else if taxable < 55800000 {
		return float64(18) / 100
	} else if taxable < 60400000 {
		return float64(19) / 100
	} else if taxable < 66700000 {
		return float64(20) / 100
	} else if taxable < 74500000 {
		return float64(21) / 100
	} else if taxable < 83200000 {
		return float64(22) / 100
	} else if taxable < 95000000 {
		return float64(23) / 100
	} else if taxable < 110000000 {
		return float64(24) / 100
	} else if taxable < 134000000 {
		return float64(25) / 100
	} else if taxable < 169000000 {
		return float64(26) / 100
	} else if taxable < 221000000 {
		return float64(27) / 100
	} else if taxable < 390000000 {
		return float64(28) / 100
	} else if taxable < 463000000 {
		return float64(39) / 100
	} else if taxable < 561000000 {
		return float64(30) / 100
	} else if taxable < 709000000 {
		return float64(31) / 100
	} else if taxable < 965000000 {
		return float64(32) / 100
	} else if taxable < 1419000000 {
		return float64(33) / 100
	} else {
		return float64(34) / 100
	}
}
