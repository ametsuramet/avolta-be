package resp

type SaleResponse struct {
	ID              string  `json:"id"`
	Date            string  `json:"date"`
	Code            string  `json:"code"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name" `
	ProductSKU      string  `json:"product_sku" `
	ShopID          string  `json:"shop_id"`
	ShopName        string  `json:"shop_name" `
	Qty             float64 `json:"qty"`
	Price           float64 `json:"price"`
	SubTotal        float64 `json:"sub_total"`
	Discount        float64 `json:"discount"`
	DiscountAmount  float64 `json:"discount_amount"`
	Total           float64 `json:"total"`
	EmployeeID      string  `json:"employee_id"`
	EmployeeName    string  `json:"employee_name"`
	EmployeePicture string  `json:"employee_picture"`
	IncentiveID     *string `json:"incentive_id"`
}
