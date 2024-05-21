package resp

type PayRollReponse struct {
	ID           string `json:"id"`
	Notes        string `json:"notes"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
}
