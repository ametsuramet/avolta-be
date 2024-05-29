package resp

type PayRollReportItemResponse struct {
	ID                     string          `json:"id"`
	TotalTakeHomePay       float64         `json:"total_take_home_pay"`
	TotalReimbursement     float64         `json:"total_reimbursement"`
	EmployeeName           string          `json:"employee_name"`
	EmployeeID             string          `json:"employee_id"`
	EmployeePhone          string          `json:"employee_phone"`
	EmployeeEmail          string          `json:"employee_email"`
	EmployeeIdentityNumber string          `json:"employee_identity_number"`
	BankAccountNumber      string          `json:"bank_account_number"`
	BankName               string          `json:"bank_name"`
	BankCode               string          `json:"bank_code"`
	BankID                 string          `json:"bank_id"`
	PayRoll                PayRollResponse `json:"pay_roll"`
	Status                 string          `json:"status" gorm:"type:enum('PROCESSING', 'FINISHED', 'PAID');default:'PROCESSING'"`
}
