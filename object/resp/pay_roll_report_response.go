package resp

import "time"

type PayRollReportResponse struct {
	ID                    string                      `json:"id"`
	Description           string                      `json:"description"`
	ReportNumber          string                      `json:"report_number"`
	UserID                string                      `json:"user_id"`
	UserName              string                      `json:"user_name"`
	StartDate             time.Time                   `json:"start_date"`
	EndDate               time.Time                   `json:"end_date"`
	Status                string                      `json:"status"`
	Items                 []PayRollReportItemResponse `json:"items"`
	GrandTotalTakeHomePay float64                     `json:"grand_total_take_home_pay"`
}
