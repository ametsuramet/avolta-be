package resp

type PayRollCostResponse struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	BpjsTkJht   bool    `json:"bpjs_tk_jht"`
	BpjsTkJp    bool    `json:"bpjs_tk_jp"`
	Tariff      float64 `json:"tariff"`
	DebtDeposit bool    `json:"debt_deposit"`
}
