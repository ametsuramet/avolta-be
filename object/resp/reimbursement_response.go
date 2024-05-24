package resp

type ReimbursementReponse struct {
	ID           string                     `json:"id"`
	Date         string                     `json:"date"`
	Name         string                     `json:"name"`
	Notes        string                     `json:"notes"`
	Total        float64                    `json:"total"`
	Balance      float64                    `json:"balance"`
	Status       string                     `json:"status" `
	Items        []ReimbursementItemReponse `json:"items"`
	EmployeeID   string                     `json:"employee_id"`
	EmployeeName string                     `json:"employee_name"`
	Transactions []TransactionReponse       `json:"transactions"`
	Attachments  []string                   `json:"attachments"`
	Remarks      string                     `json:"remarks"`
}
