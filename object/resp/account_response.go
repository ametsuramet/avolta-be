package resp

type AccountReponse struct {
	ID                    string               `json:"id"`
	Name                  string               `json:"name"`
	Code                  string               `json:"code"`
	Color                 string               `json:"color"`
	Description           string               `json:"description"`
	IsDeletable           bool                 `json:"is_deletable"`
	IsReport              bool                 `json:"is_report"`
	IsAccountReport       bool                 `json:"is_account_report"`
	IsCashflowReport      bool                 `json:"is_cashflow_report"`
	IsPDF                 bool                 `json:"is_pdf"`
	Type                  string               `json:"type"`
	Category              string               `json:"category"`
	CashflowGroup         string               `json:"cashflow_group"`
	CashflowSubGroup      string               `json:"cashflow_subgroup"`
	Transactions          []TransactionReponse `json:"transactions"`
	IsTax                 bool                 `json:"is_tax"`
	TypeLabel             string               `json:"type_label"`
	CashflowGroupLabel    string               `json:"cashflow_group_label"`
	CashflowSubGroupLabel string               `json:"cashflow_subgroup_label"`
}
