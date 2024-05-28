package resp

import "time"

type SimpleEmployeeResponse struct {
	ID         string `json:"id"`
	FullName   string `json:"full_name"`
	JobTitle   string `json:"job_title"`
	PictureUrl string `json:"picture_url"`
	UserID     string `json:"user_id"`
}
type EmployeeResponse struct {
	ID                        string             `json:"id"`
	Email                     string             `json:"email"`
	FirstName                 string             `json:"first_name"`
	MiddleName                string             `json:"middle_name"`
	LastName                  string             `json:"last_name"`
	Username                  string             `json:"username"`
	Phone                     string             `json:"phone"`
	JobTitle                  string             `json:"job_title"`
	JobTitleID                string             `json:"job_title_id"`
	UserID                    string             `json:"user_id"`
	Grade                     string             `json:"grade"`
	Address                   string             `json:"address"`
	Picture                   *string            `json:"picture"`
	PictureUrl                string             `json:"picture_url"`
	Cover                     string             `json:"cover"`
	DateOfBirth               *time.Time         `json:"date_of_birth"`
	EmployeeIdentityNumber    string             `json:"employee_identity_number"`
	EmployeeCode              string             `json:"employee_code"`
	FullName                  string             `json:"full_name"`
	BasicSalary               float64            `json:"basic_salary"`
	PositionalAllowance       float64            `json:"positional_allowance"`
	TransportAllowance        float64            `json:"transport_allowance"`
	MealAllowance             float64            `json:"meal_allowance"`
	NonTaxableIncomeLevelCode string             `json:"non_taxable_income_level_code"`
	TaxPayerNumber            string             `json:"tax_payer_number"`
	Gender                    string             `json:"gender"`
	OrganizationName          string             `json:"organization_name"`
	StartedWork               *time.Time         `json:"started_work"`
	Schedules                 []ScheduleResponse `json:"schedules"`
	TotalWorkingDays          int32              `json:"total_working_days"`
	TotalWorkingHours         float64            `json:"total_working_hours"`
	DailyWorkingHours         float64            `json:"daily_working_hours"`
	BankAccountNumber         string             `json:"bank_account_number"`
	BankID                    *string            `json:"bank_id"`
	BankName                  string             `json:"bank_name"`
}
