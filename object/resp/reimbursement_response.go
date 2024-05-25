package resp

type ReimbursementResponse struct {
	ID           string                      `json:"id"`
	Date         string                      `json:"date"`
	Name         string                      `json:"name"`
	Notes        string                      `json:"notes"`
	Total        float64                     `json:"total"`
	Balance      float64                     `json:"balance"`
	Status       string                      `json:"status" `
	Items        []ReimbursementItemResponse `json:"items"`
	EmployeeID   string                      `json:"employee_id"`
	EmployeeName string                      `json:"employee_name"`
	Transactions []TransactionResponse       `json:"transactions"`
	Attachments  []string                    `json:"attachments"`
	Remarks      string                      `json:"remarks"`
}
