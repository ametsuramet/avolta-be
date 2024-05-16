package resp

import "time"

type EmployeeReponse struct {
	ID                        string     `json:"id"`
	Email                     string     `json:"email"`
	FirstName                 string     `json:"first_name"`
	MiddleName                string     `json:"middle_name"`
	LastName                  string     `json:"last_name"`
	Username                  string     `json:"username"`
	Phone                     string     `json:"phone"`
	JobTitle                  string     `json:"job_title"`
	Grade                     string     `json:"grade"`
	Address                   string     `json:"address"`
	Picture                   *string    `json:"picture"`
	Cover                     string     `json:"cover"`
	DateOfBirth               *time.Time `json:"date_of_birth"`
	EmployeeIdentityNumber    string     `json:"employee_identity_number"`
	FullName                  string     `json:"full_name"`
	BasicSalary               float64    `json:"basic_salary"`
	PositionalAllowance       float64    `json:"positional_allowance"`
	TransportAllowance        float64    `json:"transport_allowance"`
	MealAllowance             float64    `json:"meal_allowance"`
	NonTaxableIncomeLevelCode string     `json:"non_taxable_income_level_code"`
	TaxPayerNumber            string     `json:"tax_payer_number"`
	Gender                    string     `json:"gender"`
	OrganizationName          string     `json:"organization_name"`
	StartedWork               *time.Time `json:"started_work"`
}
