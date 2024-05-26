package resp

type IncentiveSettingResponse struct {
	ID                     string  `json:"id"`
	ShopID                 string  `json:"shop_id"`
	ShopName               string  `json:"shop_name"`
	ProductCategoryID      string  `json:"product_category_id"`
	ProductCategoryName    string  `json:"product_category_name"`
	MinimumSalesTarget     float64 `json:"minimum_sales_target"`
	MaximumSalesTarget     float64 `json:"maximum_sales_target"`
	MinimumSalesCommission float64 `json:"minimum_sales_commission"`
	MaximumSalesCommission float64 `json:"maximum_sales_commission"`
	SickLeaveThreshold     float64 `json:"sick_leave_threshold"`
	OtherLeaveThreshold    float64 `json:"other_leave_threshold"`
	AbsentThreshold        float64 `json:"absent_threshold"`
}
