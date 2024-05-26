package resp

import "time"

type IncentiveReportResponse struct {
	ID           string              `json:"id"`
	Description  string              `json:"description"`
	ReportNumber string              `json:"report_number"`
	UserID       string              `json:"user_id"`
	UserName     string              `json:"user_name"`
	StartDate    time.Time           `json:"start_date"`
	EndDate      time.Time           `json:"end_date"`
	Status       string              `json:"status"`
	Incentives   []IncentiveResponse `json:"incentives"`
	Shops        []ShopResponse      `json:"shops"`
}
