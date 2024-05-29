package resp

type SettingResponse struct {
	ID                                     string  `json:"id"`
	PayRollAutoNumber                      bool    `json:"pay_roll_auto_number"`
	PayRollAutoFormat                      string  `json:"pay_roll_auto_format"`
	PayRollStaticCharacter                 string  `json:"pay_roll_static_character"`
	PayRollAutoNumberCharacterLength       int     `json:"pay_roll_auto_number_character_length"`
	PayRollPayableAccountID                *string `json:"pay_roll_payable_account_id"`
	PayRollExpenseAccountID                *string `json:"pay_roll_expense_account_id"`
	PayRollAssetAccountID                  *string `json:"pay_roll_asset_account_id"`
	PayRollTaxAccountID                    *string `json:"pay_roll_tax_account_id"`
	PayRollCostAccountID                   *string `json:"pay_roll_cost_account_id"`
	IsEffectiveRateAverage                 bool    `json:"is_effective_rate_average"`
	IsGrossUp                              bool    `json:"is_gross_up"`
	BpjsKes                                bool    `json:"bpjs_kes"`
	BpjsTkJht                              bool    `json:"bpjs_tk_jht"`
	BpjsTkJkm                              bool    `json:"bpjs_tk_jkm"`
	BpjsTkJp                               bool    `json:"bpjs_tk_jp"`
	BpjsTkJkk                              bool    `json:"bpjs_tk_jkk"`
	ReimbursementPayableAccountID          *string `json:"reimbursement_payable_account_id"`
	ReimbursementExpenseAccountID          *string `json:"reimbursement_expense_account_id"`
	ReimbursementAssetAccountID            *string `json:"reimbursement_asset_account_id"`
	IncentiveAutoNumber                    bool    `json:"incentive_auto_number"`
	IncentiveAutoFormat                    string  `json:"incentive_auto_format"`
	IncentiveStaticCharacter               string  `json:"incentive_static_character"`
	IncentiveAutoNumberCharacterLength     int     `json:"incentive_auto_number_character_length"`
	IncentiveSickLeaveThreshold            float64 `json:"incentive_sick_leave_threshold"`
	IncentiveOtherLeaveThreshold           float64 `json:"incentive_other_leave_threshold"`
	IncentiveAbsentThreshold               float64 `json:"incentive_absent_threshold"`
	PayRollReportAutoNumber                bool    `json:"pay_roll_report_auto_number"`
	PayRollReportAutoFormat                string  `json:"pay_roll_report_auto_format"`
	PayRollReportStaticCharacter           string  `json:"pay_roll_report_static_character"`
	PayRollReportAutoNumberCharacterLength int     `json:"pay_roll_report_auto_number_character_length"`
}
