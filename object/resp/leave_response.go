package resp

type LeaveResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	RequestType     string `json:"request_type"`
	LeaveCategoryID string `json:"leave_category_id"`
	LeaveCategory   string `json:"leave_category"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	EmployeeID      string `json:"employee_id"`
	EmployeeName    string `json:"employee_name"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	Remarks         string `json:"remarks"`
	AttachmentURL   string `json:"attachment_url"`
	EmployeePicture string `json:"employee_picture"`
}
