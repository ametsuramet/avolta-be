package resp

type IncentiveResponse struct {
	ID                 string                  `json:"id"`
	IncentiveReportID  string                  `json:"incentive_report_id"`
	EmployeeID         string                  `json:"employee_id"`
	EmployeeName       string                  `json:"employee_name"`
	TotalSales         float64                 `json:"total_sales"`
	TotalIncludedSales float64                 `json:"total_included_sales"`
	TotalExcludedSales float64                 `json:"total_excluded_sales"`
	TotalIncentive     float64                 `json:"total_incentive"`
	SickLeave          float64                 `json:"sick_leave"`
	OtherLeave         float64                 `json:"other_leave"`
	Absent             float64                 `json:"absent"`
	Sales              []SaleResponse          `json:"sales"`
	IncentiveShops     []IncentiveShopResponse `json:"incentive_shops"`
}
