package resp

type PayRollItemReponse struct {
	ID             string  `json:"id"`
	ItemType       string  `json:"item_type"`
	Title          string  `json:"title"`
	Notes          string  `json:"notes"`
	IsDefault      bool    `json:"is_default"`
	IsDeductible   bool    `json:"is_deductible"`
	IsTax          bool    `json:"is_tax"`
	TaxAutoCount   bool    `json:"tax_auto_count"`
	IsTaxCost      bool    `json:"is_tax_cost"`
	IsTaxAllowance bool    `json:"is_tax_allowance"`
	Amount         float64 `json:"amount"`
}
