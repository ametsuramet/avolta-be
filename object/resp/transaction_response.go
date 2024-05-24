package resp

type TransactionReponse struct {
	ID                     string               `json:"id"`
	Description            string               `json:"description"`
	Notes                  string               `json:"notes"`
	Credit                 float64              `json:"credit"`
	Debit                  float64              `json:"debit"`
	Amount                 float64              `json:"amount"`
	Date                   string               `json:"date"`
	IsIncome               bool                 `json:"is_income"`
	IsExpense              bool                 `json:"is_expense"`
	IsJournal              bool                 `json:"is_journal"`
	IsAccountReceivable    bool                 `json:"is_account_receivable"`
	IsAccountPayable       bool                 `json:"is_account_payable"`
	AccountSourceID        string               `json:"account_source_id"`
	AccountSourceName      string               `json:"account_source_name"`
	AccountDestinationID   string               `json:"account_destination_id"`
	AccountDestinationName string               `json:"account_destination_name"`
	EmployeeID             string               `json:"employee_id"`
	EmployeeName           string               `json:"employee_name"`
	IsPayRollPayment       bool                 `json:"is_pay_roll_payment"`
	IsReimbursementPayment bool                 `json:"is_reimbursement_payment"`
	TransactionRefs        []TransactionReponse `json:"transaction_refs"`
}
