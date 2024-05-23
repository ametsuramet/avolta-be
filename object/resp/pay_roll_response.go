package resp

type PayRollReponse struct {
	ID                              string               `json:"id"`
	PayRollNumber                   string               `json:"pay_roll_number"`
	Notes                           string               `json:"notes"`
	StartDate                       string               `json:"start_date"`
	EndDate                         string               `json:"end_date"`
	EmployeeID                      string               `json:"employee_id"`
	EmployeeName                    string               `json:"employee_name"`
	Items                           []PayRollItemReponse `json:"items"`
	Costs                           []PayRollCostReponse `json:"costs"`
	Transactions                    []TransactionReponse `json:"transactions"`
	PayableTransactions             []TransactionReponse `json:"payable_transactions"`
	TaxSummary                      CountTaxSummary      `gorm:"-" json:"tax_summary"`
	IsGrossUp                       bool                 `json:"is_gross_up"`
	IsEffectiveRateAverage          bool                 `json:"is_effective_rate_average"`
	TotalIncome                     float64              `json:"total_income"`
	TotalReimbursement              float64              `json:"total_reimbursement"`
	TotalDeduction                  float64              `json:"total_deduction"`
	TotalTax                        float64              `json:"total_tax"`
	TaxCost                         float64              `json:"tax_cost"`
	NetIncome                       float64              `json:"net_income"`
	NetIncomeBeforeTaxCost          float64              `json:"net_income_before_tax_cost"`
	TakeHomePay                     float64              `json:"take_home_pay"`
	TotalPayable                    float64              `json:"total_payable"`
	TaxAllowance                    float64              `json:"tax_allowance"`
	TakeHomePayCounted              string               `json:"take_home_pay_counted" `
	TakeHomePayReimbursementCounted string               `json:"take_home_pay_reimbursement_counted" `
	Status                          string               `json:"status" `
}

type CountTaxSummary struct {
	JobExpenseMonthly               float64 `json:"job_expense_monthly"`
	JobExpenseYearly                float64 `json:"job_expense_yearly"`
	PtkpYearly                      float64 `json:"ptkp_yearly"`
	GrossIncomeMonthly              float64 `json:"gross_income_monthly"`
	GrossIncomeYearly               float64 `json:"gross_income_yearly"`
	PkpMonthly                      float64 `json:"pkp_monthly"`
	PkpYearly                       float64 `json:"pkp_yearly"`
	TaxYearlyBasedOnProgressiveRate float64 `json:"tax_yearly_based_on_progressive_rate"`
	TaxYearly                       float64 `json:"tax_yearly"`
	TaxMonthly                      float64 `json:"tax_monthly"`
	NetIncomeMonthly                float64 `json:"net_income_monthly"`
	NetIncomeYearly                 float64 `json:"net_income_yearly"`
	CutoffPensiunMonthly            float64 `json:"cutoff_pensiun_monthly"`
	CutoffPensiunYearly             float64 `json:"cutoff_pensiun_yearly"`
	CutoffMonthly                   float64 `json:"cutoff_monthly"`
	CutoffYearly                    float64 `json:"cutoff_yearly"`
	Ter                             float64 `json:"ter"`
}
