package resp

type IncentiveShopResponse struct {
	ID                 string  `json:"id"`
	ShopID             string  `json:"shop_id"`
	ShopName           string  `json:"shop_name"`
	IncentiveID        string  `json:"incentive_id"`
	TotalSales         float64 `json:"total_sales"`
	TotalIncludedSales float64 `json:"total_included_sales"`
	TotalExcludedSales float64 `json:"total_excluded_sales"`
	TotalIncentive     float64 `json:"total_incentive"`
}
