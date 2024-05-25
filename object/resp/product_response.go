package resp

type ProductResponse struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	SKU                 string  `json:"sku"`
	Barcode             string  `json:"barcode"`
	SellingPrice        float64 `json:"selling_price"`
	ProductCategoryID   string  `json:"product_category_id"`
	ProductCategoryName string  `json:"product_category_name"`
}
